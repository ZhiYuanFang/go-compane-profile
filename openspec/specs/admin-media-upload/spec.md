# admin-media-upload

## Purpose

admin-media-upload 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/admin-deferred-oss-upload/specs/admin-deferred-upload, changes/admin-eager-image-upload/specs/admin-eager-upload-queue, changes/admin-eager-image-upload/specs/admin-image-size-3m, changes/admin-gallery-preview-fix-concurrency-1/specs/admin-upload-concurrency-1, changes/admin-image-upload-parallel/specs/admin-image-size-align, changes/admin-image-upload-parallel/specs/admin-upload-concurrency, changes/admin-oss-upload-progress-stream/specs/admin-oss-progress-stream, changes/admin-upload-pipeline-ab/specs/admin-compress-upload-decouple -->

## Requirements

### Requirement: Compress on select without OSS upload

Across all admin dual-image controls, when the operator selects or drops an image file, the system SHALL show a local preview and compress the file into pending original and thumb blobs, and MUST NOT call the OSS upload API until the enclosing form save begins.

#### Scenario: Select image on portfolio gallery

- **WHEN** the operator selects an image for a portfolio gallery slot
- **THEN** the slot shows a local preview and reaches a ready-with-pending state without a network upload request

#### Scenario: Select image on company logo

- **WHEN** the operator selects an image for the company logo field
- **THEN** the system compresses locally into pending blobs and does not upload to OSS until company save

### Requirement: Upload pending images on save with progress

When the operator saves an admin form that contains dual-image slots with pending blobs, the system SHALL upload those pending slots via the existing streaming upload path, SHALL show per-slot upload progress, and SHALL persist CMS payloads only after all required uploads for that save succeed.

#### Scenario: Save portfolio with new pending cover and gallery images

- **WHEN** the operator clicks save and one or more slots have pending blobs
- **THEN** those slots upload with progress UI and the create/update APIs are called only after every pending upload for that save completes successfully

#### Scenario: Upload failure blocks persistence

- **WHEN** any pending upload fails during save
- **THEN** the system does not complete the CMS write for that submit with missing image URLs, and the failed slot shows an error state for retry

### Requirement: Unsaved removal creates no OSS objects

When the operator removes or clears a dual-image slot that has never been persisted with remote URLs for the current edit (pending-only or still uploading locally), the system SHALL discard local preview and pending blobs and MUST NOT issue an OSS delete or upload request for that removal.

#### Scenario: Delete gallery image before save

- **WHEN** the operator uploads/selects a gallery image (compress completes) and then removes that slot without saving
- **THEN** no OSS upload or delete network request is triggered for that removal

### Requirement: Ready pending does not block save initiation

Slots that are compressed and waiting with pending blobs (not currently compressing or uploading, and not in error) MUST allow the operator to initiate save; the save flow itself performs the deferred uploads.

#### Scenario: Save enabled when images are ready locally

- **WHEN** all non-empty dual slots are either already remote-DONE or ready-with-pending, and none are compressing, uploading, or in error
- **THEN** the save control is not blocked by a “still uploading” reason solely due to pending blobs

### Requirement: Site-wide deferred behavior

Deferred upload SHALL apply consistently to every admin surface that uses dual-image upload (including company, portfolio cover and gallery, activity, and pricing), not only portfolio gallery.

#### Scenario: Activity image follows same pipeline

- **WHEN** the operator selects an activity cover image and later saves the activity
- **THEN** OSS upload occurs during save with progress, not immediately on select

### Requirement: Delete cancels queued or in-flight work
Removing a gallery slot (single remove or batch delete) SHALL cancel that slot's queued job or abort an in-flight compress/upload when possible. A late successful upload response for a removed `slotId` MUST NOT re-add the image to the list.

#### Scenario: Delete while queued
- **WHEN** an admin deletes a slot that is still queued
- **THEN** that job never uploads

#### Scenario: Delete while uploading
- **WHEN** an admin deletes a slot that is uploading
- **THEN** the client aborts or ignores the result for that slotId
- **AND** the slot does not reappear after a late response

### Requirement: Save refuses unsettled slots
Save SHALL be blocked with a clear reason while any non-empty slot is still compressing, queued for compression, or uploading, and SHALL remain blocked while any slot is in a failed state. Slots that already have remote URLs SHALL be submitted as URLs without re-uploading.

#### Scenario: Block save on failure
- **WHEN** any slot shows a failed upload status
- **THEN** save does not complete successfully until the admin retries successfully or removes the slot

#### Scenario: Block while processing
- **WHEN** any slot is still compressing, queued, or uploading
- **THEN** save is blocked with clear feedback until that work finishes or fails

### Requirement: Original max 3MiB
Client compression and server upload validation SHALL enforce an original image maximum of 3 MiB (`3 * 1024 * 1024` bytes).

#### Scenario: Client compress target
- **WHEN** a selected original exceeds 3MiB
- **THEN** the client attempts to compress it under 3MiB before upload

#### Scenario: Server rejects oversized original
- **WHEN** an upload presents an original larger than 3MiB
- **THEN** the server rejects it with an invalid-parameter style error

### Requirement: Thumb max aligned at 0.5MiB
Client and server SHALL both treat the thumbnail maximum as `0.5 * 1024 * 1024` bytes so a thumb accepted by client compression is not rejected solely for exceeding `500 * 1024`.

#### Scenario: Client and server agree
- **WHEN** a thumb size is between `500*1024` and `0.5*1024*1024` bytes and within the compress success bound
- **THEN** both client acceptance and server validation allow it

### Requirement: Compress concurrency is two and upload concurrency is one
管理端共享压缩池与上传池 MUST 使用相互独立的并发上限：压缩/水印池最多同时运行 2 个任务（`COMPRESS_CONCURRENCY = 2`），上传池最多同时进行 1 个 HTTP 上传（`UPLOAD_CONCURRENCY = 1`）。压缩任务 MUST NOT 占用上传池槽位。

#### Scenario: Second upload waits for first
- **WHEN** two slots are ready to upload
- **THEN** only one HTTP upload request runs at a time
- **AND** the second starts only after the first completes, fails, or is cancelled

#### Scenario: Second compress waits for first
- **WHEN** more than two slots need compression
- **THEN** at most two compressions run at a time
- **AND** the next compression starts only after one completes, fails, or is cancelled
- **AND** slots waiting for compress MUST still show their local previews if set

#### Scenario: Two compressions run together
- **WHEN** exactly two slots need compression and none is in flight
- **THEN** both compressions run concurrently

### Requirement: Aligned thumbnail size limit
Client-side image processing SHALL use the same byte ceiling for thumbnail success checks as the compression target of 0.5 MiB (`0.5 * 1024 * 1024` bytes). A thumbnail that successfully finishes compression under that target MUST NOT be rejected solely for exceeding `500 * 1024` bytes.

#### Scenario: Thumb in former dead zone succeeds
- **WHEN** compression produces a thumbnail whose size is between `500*1024` and `0.5*1024*1024` bytes inclusive of the compress success bound
- **THEN** `processDualImage` accepts it as a valid thumb
- **AND** does not throw「缩略图超过 500KB…」for that reason alone

### Requirement: Consistent limit constants
Original and thumb max sizes SHALL be derived from the same unit convention as `compressUnder`'s `maxSizeMB * 1024 * 1024` checks so post-conditions cannot disagree with compress success.

#### Scenario: Single source of truth
- **WHEN** thumb max is defined for processing
- **THEN** it equals `0.5 * 1024 * 1024` (or an equivalent expression of 0.5 MiB)

### Requirement: Pending uploads share the serial upload queue
When resolving multiple pending dual images for a save, the admin client SHALL enqueue them on the shared upload queue with a maximum concurrency of 1, regardless of the compress pool concurrency. Items that already have remote URLs and no pending blobs SHALL not be uploaded again.

#### Scenario: Multiple pending uploads are serialized
- **WHEN** an admin saves a portfolio with multiple pending gallery images
- **THEN** at most one upload request is in flight at a time
- **AND** all pending images are uploaded before the portfolio/gallery write continues

#### Scenario: Preserve gallery order
- **WHEN** pending images at positions 1..N are resolved via the pool
- **THEN** the resolved URL list retains the same order as the input gallery list

### Requirement: Upload progress during save
While save-time uploads are running, the admin UI SHALL show per-slot upload progress from the streaming upload events, and the form save control SHALL indicate that upload and save are in progress.

#### Scenario: Per-slot progress updates
- **WHEN** pending dual uploads are in progress during save
- **THEN** each uploading slot shows its own stream percent (e.g.「上传中 42%」)
- **AND** the save control shows an uploading-and-saving state

### Requirement: Fail save on any upload error
If any pending dual upload fails, the save operation SHALL fail and MUST NOT treat the gallery as successfully updated for that submit.

#### Scenario: One upload fails
- **WHEN** one of several parallel uploads returns an error
- **THEN** the save surfaces an error
- **AND** does not complete as a successful save

### Requirement: Streaming upload with OSS progress events
The system SHALL provide an authenticated admin upload endpoint that accepts the existing dual-image multipart fields (`original`, `thumb`, `category`), uploads both objects to OSS, and **streams progress events** to the client while OSS puts are in progress—not only after both puts complete.

#### Scenario: Progress events before done
- **WHEN** an admin client uploads a dual image pair via the streaming endpoint
- **THEN** the client receives one or more progress events with an overall percent before the final success event
- **AND** the final success event includes CDN URLs for original and thumb

#### Scenario: Error event on failure
- **WHEN** OSS upload fails after the multipart body was accepted
- **THEN** the stream ends with an error event containing a message
- **AND** the client MUST NOT treat the upload as successful

### Requirement: Progress reflects OSS body read
Progress percent SHALL be derived primarily from bytes consumed while uploading object bodies to OSS (original and thumb), weighted by their known sizes, so that progress continues to advance after the browser has finished sending the HTTP request body.

#### Scenario: Percent advances during OSS put
- **WHEN** the server is putting the original object to OSS
- **THEN** progress percent increases as the original body is read for upload
- **AND** thumb upload continues to advance overall percent until completion

### Requirement: Admin UI shows stream progress
The admin upload pipeline SHALL use the streaming endpoint and display the stream’s overall percent on the slot status (e.g.「上传中 N%」) until success or failure. Cancel/abort SHALL abort the in-flight stream request.

#### Scenario: Badge tracks OSS-weighted percent
- **WHEN** the stream reports progress pct=42
- **THEN** the gallery/cover status shows uploading progress at 42%

#### Scenario: Done clears progress
- **WHEN** the stream reports done with URLs
- **THEN** the slot stores CDN URLs, status becomes uploaded, and progress is cleared

### Requirement: Compress outside upload queue
Gallery (and cover when using the shared pipeline) image compression SHALL NOT occupy a slot in the shared upload concurrency pool. Compression SHALL run in a separate pool with a maximum concurrency of 2; only the HTTP upload of already-compressed pending blobs SHALL use the upload queue with a maximum concurrency of 1.

#### Scenario: Compress while an upload is in progress
- **WHEN** a slot's compression is pending while another slot's upload is in progress
- **THEN** the compression MAY begin or complete without waiting for an upload pool slot
- **AND** at most one HTTP upload runs at a time

#### Scenario: Upload queue receives pending blobs only
- **WHEN** a save flushes a slot whose compression finished successfully
- **THEN** the client enqueues an upload of the compressed original+thumb pair
- **AND** the upload job does not re-run full dual compression

### Requirement: Parallel compress with bounded pool
The client SHALL compress multiple pending files with a maximum concurrency of 2 in the shared compress pool (`COMPRESS_CONCURRENCY = 2`). Slots waiting for a compress slot SHALL keep their local preview and a waiting/processing status as appropriate.

#### Scenario: Compress pool capped
- **WHEN** multiple slots need compression at once
- **THEN** at most two compressions run at a time
- **AND** remaining slots wait without blocking UI display of their previews

### Requirement: Cancel aborts compress and upload
Deleting or clearing a slot SHALL cancel that slot’s compress work and upload work. A late compress or upload completion for a cancelled or removed `slotId` MUST NOT restore the slot into the list.

#### Scenario: Delete during compress
- **WHEN** an admin deletes a slot that is compressing
- **THEN** compression for that slot is cancelled or its result discarded
- **AND** no upload is started for that slotId

### Requirement: 选图压缩不自动加水印
管理端为内容图（作品封面、作品图、活动封面）选择或拖入图片并完成压缩后，客户端 MUST NOT 自动叠加水印；槽位 SHALL 以未加水印的 pending original/thumb 进入待保存，并允许直接保存上线。公司 logo MUST NOT 加水印。

#### Scenario: 作品图选图后无水印
- **WHEN** 运营为作品图槽位选择图片并完成压缩
- **THEN** 该槽位进入待保存，其 pending 的 original 与 thumb 不含 `MUHU` 水印

#### Scenario: 活动封面选图后无水印
- **WHEN** 运营为活动封面选择图片并完成压缩
- **THEN** 该槽位进入待保存，pending 不含水印

#### Scenario: 干净图可直接保存上线
- **WHEN** 运营保存仅含压缩后、未加水印 pending 的内容图表单
- **THEN** 上传到服务端的 original 与 thumb 为未加水印的压缩结果
- **AND** 服务端返回的 CDN URL 指向无水印图

#### Scenario: 公司 logo 仍不加水印
- **WHEN** 运营为公司 Logo 选择图片并完成压缩
- **THEN** 该图片不叠加水印，字节与压缩结果一致

### Requirement: 显式加水印（已上线或干净 pending）
运营通过「加水印」控件对内容图槽位请求加水印时，客户端 SHALL 为该槽位叠加 `MUHU` 水印并写入待上传的 original 与 thumb。目标槽位 SHALL 包括：已上线且无 pending 的内容图，以及已有干净 pending（待保存、尚未上线）的内容图。公司 logo MUST NOT 作为加水印目标。重复加水印 MAY 叠加。

#### Scenario: 已上线内容图加水印
- **WHEN** 运营对已上线且无 pending 的内容图点击加水印
- **THEN** 客户端取得原图像素后进入压缩与水印管线
- **AND** 槽位进入待保存，pending 的 original 与 thumb 含 `MUHU` 水印

#### Scenario: 干净 pending 加水印
- **WHEN** 运营对处于待保存、pending 尚未上线的内容图点击加水印
- **THEN** 客户端直接对本地 pending 字节叠加水印，不发起 CDN 下载
- **AND** 槽位仍为待保存，pending 更新为含水印的 original 与 thumb

#### Scenario: logo 不可加水印
- **WHEN** 槽位类别为公司 Logo
- **THEN** 「加水印」控件不可用或对 logo 无效

#### Scenario: 保存上传显式水印后的字节
- **WHEN** 运营保存经显式加水印后进入待保存的内容图
- **THEN** 上传到服务端的 original 与 thumb 均为加水印后的图片
- **AND** 服务端返回的 CDN URL 指向水印图

### Requirement: 水印斜向重复铺满
水印 SHALL 以倾斜角度重复绘制 `MUHU` 文字并铺满整张图片，密度保持稀疏且覆盖图片四角；字号与步距 MUST 按图片尺寸等比缩放，使 original 与 thumb 呈现一致的视觉密度。

#### Scenario: 多处重复且覆盖四角
- **WHEN** 客户端为一张内容图生成水印
- **THEN** 图中可见多处 `MUHU` 文字
- **AND** 图片四个角落附近均有水印覆盖，裁掉单边无法去除水印

#### Scenario: 稀疏不遮挡主体
- **WHEN** 查看一张加水印后的 16:9 内容图
- **THEN** 可见的重复文字条数在 4 到 16 之间
- **AND** 图片主体区域仍然清晰可辨

#### Scenario: 原图与缩略图密度一致
- **WHEN** 同一张图分别生成 original 与 thumb
- **THEN** 两张图的文字相对尺寸与重复密度基本一致

### Requirement: 水印后保持上传体积上限
水印重编码后，original MUST 不超过 3MiB、thumb MUST 不超过 0.5MiB。超过上限时，客户端 SHALL 先降低编码质量、必要时缩小尺寸后重新编码。若降低质量并缩小尺寸后仍不满足上限，槽位 MUST 进入失败状态。

#### Scenario: JPEG 超限降质重编码
- **WHEN** 水印后的 JPEG 超过对应上限
- **THEN** 客户端以更低质量重新编码，直到不超过上限或达到最低质量

#### Scenario: 仍超限则失败
- **WHEN** 降质与缩尺寸后仍无法满足上限
- **THEN** 槽位显示处理失败，且不进入待保存状态

### Requirement: 水印失败不得静默降级
在运营主动请求加水印的路径上，水印绘制或重编码失败时，槽位 SHALL 进入失败状态并向运营展示错误；客户端 MUST NOT 以未加水印的图片冒充本次加水印成功而进入待保存或上传。选图/压缩未请求加水印时，本要求不适用。

#### Scenario: 绘制失败
- **WHEN** 运营点击加水印后客户端无法完成水印绘制或编码
- **THEN** 槽位状态为失败并显示错误
- **AND** 不得将本次加水印视为成功

#### Scenario: 选图压缩失败与水印无关
- **WHEN** 运营仅选择图片且压缩失败
- **THEN** 槽位按既有压缩失败路径处理
- **AND** 不因「未加水印」单独判失败

### Requirement: 补水印原图浏览器直拉 CDN
为已上线内容图补加水印时，客户端 SHALL 使用浏览器对槽位已有 CDN URL（优先 original，其次 thumb）发起跨域 `fetch` 取得图片字节，再进入既有压缩/水印管线。系统 MUST NOT 再通过同源 admin 媒体代理端点代下这些字节。本机或其它未获 CDN CORS 放行的 Origin 直拉失败时，客户端 MUST 按现有下载失败路径提示错误，且 MUST NOT 回退到服务端代理。

#### Scenario: 生产管理端直拉 CDN 成功
- **WHEN** 运营在已配置 CORS 的管理端域名上对已上传内容图点击加水印
- **THEN** 浏览器直接向该图 CDN URL 请求图片字节
- **AND** 取得的字节进入后续压缩与水印，不请求 `/admin/api/media/proxy`

#### Scenario: 未放行 Origin 失败且不回退代理
- **WHEN** 当前页面 Origin 未被 CDN CORS 允许，导致直拉失败
- **THEN** 该槽位按下载失败处理并向运营展示错误
- **AND** 客户端不调用 admin 媒体代理作为兜底
