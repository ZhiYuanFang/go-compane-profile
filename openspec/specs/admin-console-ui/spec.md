# admin-console-ui

## Purpose

admin-console-ui 域行为规格（自变更 delta / 旧基线迁移）。

<!-- migrated-from: changes/admin-eager-image-upload/specs/admin-image-slot-status, changes/admin-gallery-grid-batch/specs/admin-gallery-batch-delete, changes/admin-gallery-grid-batch/specs/admin-gallery-batch-import, changes/admin-gallery-grid-batch/specs/admin-gallery-grid, changes/admin-gallery-preview-fix-concurrency-1/specs/admin-gallery-batch-preview, changes/admin-image-drag-drop/specs/admin-image-drag-drop, changes/admin-portfolio-gallery-ux/specs/admin-image-preview, changes/admin-portfolio-gallery-ux/specs/admin-portfolio-pair-gallery, changes/admin-slot-status-colors/specs/admin-slot-status-tones, changes/admin-upload-pipeline-ab/specs/admin-gallery-instant-preview, changes/company-profile-api-admin/specs/admin-console-ui, changes/portfolio-optional-gallery-pair/specs/optional-gallery-pair -->

## Requirements

### Requirement: Per-slot upload status badge
Each gallery thumbnail (and cover dual field) SHALL show a Simplified Chinese status reflecting the slot's upload pipeline state.

#### Scenario: Status vocabulary
- **WHEN** a slot is compressing, queued, uploading, successfully uploaded, or failed
- **THEN** the UI shows a corresponding label such as 处理中…、排队中、上传中、已上传、失败

### Requirement: Remote URL shows uploaded
A slot that already has remote thumb/original URLs and no pending upload work SHALL display 已上传 (not 待上传).

#### Scenario: Loaded portfolio images
- **WHEN** an admin opens an existing portfolio whose gallery images already have CDN URLs
- **THEN** those slots show 已上传

### Requirement: Failed state allows retry
A failed upload status SHALL remain visible until the admin retries or removes the slot, and a retry control SHALL be available.

#### Scenario: Retry after failure
- **WHEN** an admin retries a failed slot
- **THEN** the slot re-enters the upload pipeline and status updates accordingly

### Requirement: Select gallery images with checkboxes
The gallery grid SHALL allow selecting any subset of slots via per-cell checkboxes.

#### Scenario: Toggle selection
- **WHEN** an admin checks a cell checkbox
- **THEN** that slot is included in the current selection
- **WHEN** they uncheck it
- **THEN** it is removed from the selection

### Requirement: Select all and clear selection
The gallery toolbar SHALL provide 全选 and 取消全选 (or equivalent toggle) so all current slots can be selected or deselected in one action.

#### Scenario: Select all
- **WHEN** an admin chooses 全选 with N gallery slots present
- **THEN** all N slots are selected

#### Scenario: Clear selection
- **WHEN** an admin chooses 取消全选
- **THEN** no slots remain selected

### Requirement: Delete selected with confirmation
The toolbar SHALL provide 删除所选. When one or more slots are selected, activating it SHALL ask for confirmation including the count, then remove those slots from the gallery model and clear the selection. When none are selected, 删除所选 SHALL be disabled or no-op.

#### Scenario: Bulk delete confirmed
- **WHEN** an admin selects multiple slots and confirms 删除所选
- **THEN** those slots are removed from the gallery list
- **AND** selection is cleared

#### Scenario: Bulk delete cancelled
- **WHEN** an admin cancels the delete confirmation
- **THEN** the gallery list and selection remain unchanged

#### Scenario: Delete selected disabled when empty
- **WHEN** no slots are selected
- **THEN** 删除所选 cannot remove any items

### Requirement: Simplified Chinese batch-delete copy
Admin labels for select-all, clear selection, delete selected, and confirm dialogs SHALL use Simplified Chinese.

#### Scenario: Chinese labels
- **WHEN** an admin views gallery batch-delete controls
- **THEN** they see Chinese copy such as 全选、取消全选、删除所选

### Requirement: Pending images show local preview
When a gallery dual has `pending` blobs and no remote URLs yet, the admin UI SHALL show a local object-URL preview of the pending thumb (or equivalent) instead of a blank preview area.

#### Scenario: Batch-added pending preview
- **WHEN** batch select finishes processing a file into a pending dual
- **THEN** that slot's thumbnail/preview shows the local image
- **AND** status may still indicate pending upload until form submit

### Requirement: Progressive batch select with progress
Batch file selection SHALL process images and append each completed slot to the list progressively, showing progress such as processed/total while work is in progress. The UI MUST remain understandable during long batches (not a blank list until the entire batch completes).

#### Scenario: Progressive append
- **WHEN** an admin batch-selects multiple images
- **THEN** slots appear as each file finishes processing
- **AND** a progress indicator shows how many of the batch are done

#### Scenario: Partial failure
- **WHEN** one file in a batch fails processing after some succeeded
- **THEN** successfully processed slots remain in the list
- **AND** an error message explains the failure

### Requirement: Batch controls disabled while processing
While batch processing is running, batch select and destructive batch actions SHALL be disabled or otherwise blocked to avoid conflicting edits.

#### Scenario: Processing lock
- **WHEN** batch processing is in progress
- **THEN** starting another batch select is disabled until processing finishes

### Requirement: Two-column thumbnail gallery grid
The admin portfolio gallery editor SHALL display gallery slots in a two-column grid of compact thumbnails (not full-width DualImageField rows as the primary list surface). On narrow viewports the grid MAY collapse to one column.

#### Scenario: Dense grid layout
- **WHEN** an admin opens a portfolio with multiple gallery images
- **THEN** thumbnails appear in a two-column grid
- **AND** each cell shows an index and a small preview when an image is present

### Requirement: Click thumbnail to edit
Clicking a gallery thumbnail SHALL open an edit surface (dialog or equivalent) that exposes DualImageField controls for that slot (replace, clear, full preview). Checkbox clicks SHALL NOT open the editor.

#### Scenario: Open slot editor
- **WHEN** an admin clicks a gallery thumbnail (not its checkbox)
- **THEN** an editor opens for that slot with dual-image controls
- **AND** closing the editor returns to the grid with any changes reflected

#### Scenario: Empty slot still editable
- **WHEN** an admin clicks an empty gallery thumbnail/placeholder
- **THEN** the editor opens so they can choose an image for that slot

### Requirement: Compact per-slot affordances
Each grid cell SHALL support selection via checkbox. Reorder (上移/下移) and single remove MAY live in the edit dialog and/or as compact cell actions; the primary list SHALL NOT require a full-width DualImageField per row.

#### Scenario: Cell shows selection control
- **WHEN** gallery items are shown in the grid
- **THEN** each cell has a checkbox for multi-select

### Requirement: Batch select keeps all local previews
When an admin batch-selects multiple images, the client SHALL associate a local preview with every new slot in the same update cycle (or equivalent) such that no selected slot loses its preview due to concurrent status patches on other slots.

#### Scenario: Ten files all show previews immediately
- **WHEN** an admin batch-selects ten images
- **THEN** all ten gallery cells show a local preview without waiting for compression or upload to finish
- **AND** subsequent status updates on one slot MUST NOT clear another slot’s local preview

### Requirement: Async patches do not clobber sibling slots
Per-slot pipeline patches (queued / compressing / uploading / done / error) SHALL update only the targeted `slotId` while preserving other slots’ fields including `localPreview`, `pending`, and remote URLs.

#### Scenario: Compress patch preserves neighbors
- **WHEN** slot A’s status is patched to compressing while slot B already has a localPreview
- **THEN** slot B’s localPreview remains visible

### Requirement: Drag-and-drop image selection on dual-image fields
The admin dual-image field SHALL accept an image file dropped onto the field surface. The dropped file SHALL enter the same local dual-image processing pipeline used by the file picker (pending local pair until form submit). Non-image drops SHALL be rejected with an error message and SHALL NOT change the current value.

#### Scenario: Drop image onto empty field
- **WHEN** an operator drops a valid image file onto an empty dual-image field
- **THEN** the field processes that image and shows a local preview without requiring the system file dialog

#### Scenario: Reject non-image drop
- **WHEN** an operator drops a non-image file onto a dual-image field
- **THEN** the field shows an error and keeps the previous value unchanged

### Requirement: Confirm before replacing existing preview via drop
When the dual-image field already has a preview (saved or pending) and the operator drops a new image onto it, the system SHALL ask for confirmation before replacing. If the operator cancels, the current value SHALL remain unchanged.

#### Scenario: Confirm replace on drop
- **WHEN** an operator drops an image onto a field that already shows a preview and confirms the replace prompt
- **THEN** the field processes the new image and replaces the previous pending or displayed selection

#### Scenario: Cancel replace on drop
- **WHEN** an operator drops an image onto a field that already shows a preview and cancels the replace prompt
- **THEN** the field keeps the previous preview and value

### Requirement: Empty preview click opens file picker
When the dual-image preview is empty, activating the preview (click) SHALL open the same file picker used by the choose-image control. When the preview is non-empty, activating the preview SHALL open the existing enlarge/lightbox behavior and SHALL NOT open the file picker.

#### Scenario: Click empty preview to choose file
- **WHEN** an operator clicks an empty dual-image preview
- **THEN** the system file picker opens for image selection

#### Scenario: Click filled preview enlarges
- **WHEN** an operator clicks a non-empty dual-image preview
- **THEN** the enlarged lightbox opens and the file picker does not open

### Requirement: File button remains available
The dual-image field SHALL continue to provide an explicit control to choose or re-choose an image via the system file picker, in addition to drag-and-drop and empty-preview click.

#### Scenario: Button still works
- **WHEN** an operator uses the choose/re-choose button
- **THEN** the system file picker opens as before

### Requirement: Contain fit for dual-image preview
The admin dual-image field SHALL display the selected or saved image entirely within the existing preview box using contain-fit behavior (no cropping). The preview box dimensions SHALL remain unchanged from the current layout size.

#### Scenario: Selected image fits inside preview
- **WHEN** an operator selects or views an image in a dual-image field
- **THEN** the full image is visible inside the preview box without cropping, even if letterboxing appears inside the box

### Requirement: Click to enlarge dual-image preview
The admin dual-image field SHALL allow the operator to open an enlarged view of the current image by activating the preview (for example, clicking the preview). The enlarged view SHALL prefer the original URL when available, otherwise the current preview source (including a local pending preview). The operator SHALL be able to dismiss the enlarged view.

#### Scenario: Open and dismiss enlarge
- **WHEN** an operator activates a non-empty dual-image preview
- **THEN** an enlarged view of that image is shown and can be dismissed to return to the form

#### Scenario: Empty preview not enlargeable
- **WHEN** the dual-image preview has no image
- **THEN** activating the preview does not open an enlarged view

### Requirement: Enlarge applies to all dual-image fields
Every admin screen that uses the dual-image field (company logos, portfolio cover, pair gallery items, pricing) SHALL provide the same click-to-enlarge behavior. Portfolio list cover thumbnails are out of scope for enlarge.

#### Scenario: Pricing image enlarge
- **WHEN** an operator views a selected pricing image preview and activates it
- **THEN** the enlarged view opens for that image

### Requirement: Hide portfolio slug in admin editor
The portfolio create and edit forms SHALL NOT present a slug input. On create, the client SHALL omit a custom slug (empty) so the server may auto-generate one. On edit, the client SHALL NOT require the operator to manage slug.

#### Scenario: Create without slug field
- **WHEN** an operator opens the new portfolio form
- **THEN** no slug field is shown and a successful create still persists a server-generated slug

### Requirement: Paired render and real gallery editor
The portfolio editor SHALL present gallery images as ordered pairs, each pair containing one render (effect) image field and one optional real (site photo) image field. Adding a pair SHALL add both slots together. Reordering or removing SHALL apply to the whole pair. Real MAY be empty.

#### Scenario: Add pair with optional real
- **WHEN** an operator adds a gallery pair and fills only the render image
- **THEN** the pair is valid to keep in the editor and the real slot may remain empty

#### Scenario: Reorder pairs
- **WHEN** an operator moves a pair up or down
- **THEN** both the render and real for that pair move together

### Requirement: Persist pairs as aligned render and real arrays
On save, the editor SHALL submit `renders` and `reals` arrays of equal length corresponding to pair order. An empty real SHALL be represented as an empty dual-URL placeholder at the same index so public clients can align by index.

#### Scenario: Save with a middle empty real
- **WHEN** an operator saves three pairs where the second real is empty
- **THEN** the saved `reals` list has length three and the second entry has empty thumb and original URLs

#### Scenario: Load unequal historical lists
- **WHEN** the editor loads a portfolio whose stored render and real list lengths differ
- **THEN** the UI zips them into pairs by index and pads the shorter side with empty dual slots so no images are dropped

### Requirement: Navigate to portfolio list after successful save
After a successful portfolio create or update, the admin console SHALL navigate the operator to the portfolio list page.

#### Scenario: Create returns to list
- **WHEN** an operator successfully creates a portfolio
- **THEN** the console navigates to the portfolio list

#### Scenario: Update returns to list
- **WHEN** an operator successfully updates a portfolio
- **THEN** the console navigates to the portfolio list

### Requirement: Four status tone tiers
The admin UI SHALL classify each image slot status into one of four visual tones: wait (排队中), busy (处理中 or 上传中), done (已上传), or error (失败). Compressing and uploading SHALL share the busy tone.

#### Scenario: Queued is wait
- **WHEN** a slot status is queued
- **THEN** its tone is wait

#### Scenario: Compressing and uploading are busy
- **WHEN** a slot status is compressing or uploading
- **THEN** its tone is busy

#### Scenario: Done and error tones
- **WHEN** a slot status is done
- **THEN** its tone is done
- **WHEN** a slot status is error
- **THEN** its tone is error

### Requirement: Status text color by tone
Gallery status badges and cover DualImageField status text SHALL use distinct colors per tone so wait, busy, done, and error are visually distinguishable at a glance.

#### Scenario: Badge colors differ by tone
- **WHEN** several gallery slots show wait, busy, done, and error statuses
- **THEN** their status label colors differ according to those four tones

### Requirement: Border color by tone
Gallery cells and the cover dual container SHALL apply a border (outline) color matching the same tone as the status text when a tone is present.

#### Scenario: Cell border matches status tone
- **WHEN** a gallery slot is uploading (busy)
- **THEN** the cell border uses the busy tone color
- **AND** the status text also uses the busy tone color

### Requirement: Immediate local preview on file select
When an admin selects one or more images for gallery slots (single replace, empty pick, or batch), the client SHALL show a local preview for each selected file as soon as the file is associated with a slot, without waiting for compression or upload to start or finish.

#### Scenario: Batch select shows all previews
- **WHEN** an admin batch-selects multiple images
- **THEN** each new slot displays a local preview of its file promptly
- **AND** previews appear even for slots that are still waiting to compress or upload

#### Scenario: Replace updates preview immediately
- **WHEN** an admin replaces an image on a slot
- **THEN** the cell preview updates to the newly selected file before upload completes

### Requirement: Preview available while busy
While a slot is compressing, queued, or uploading, if a local or pending preview exists, the admin SHALL be able to open that preview (e.g. lightbox). Busy state MUST NOT blank or hide an available preview.

#### Scenario: Lightbox during upload
- **WHEN** a slot is uploading and has a local or pending preview
- **THEN** the admin can open the preview viewer for that slot

### Requirement: Vue 3 admin console with glassmorphism gray theme
The system SHALL provide a Vue 3 admin SPA with a glassmorphism visual language and high-end gray theme (deep gray atmosphere, translucent panels, cool silver accents). The SPA SHALL be served from the application host under `/admin/`.

#### Scenario: Admin entry
- **WHEN** an operator opens `/admin/` on the public host
- **THEN** the Vue 3 admin application loads and presents a login gate before content management

### Requirement: Content management screens
The admin console SHALL provide screens to edit company profile, manage ordered portfolios (including cover and render/real galleries), and manage the pricing dual image.

#### Scenario: Navigate modules
- **WHEN** an authenticated operator uses the admin navigation
- **THEN** they can reach company, portfolio, and pricing management views

### Requirement: Local dual-image UX on image fields
For every image field (logos, covers, gallery items, pricing), the console SHALL show local preview of the processed pair, allow re-selection before submit, and upload pairs only on form submit together with the rest of the save operation (or explicitly as part of the save flow such that unused drafts are not persisted to OSS).

#### Scenario: Preview before submit
- **WHEN** the operator selects an image for a field
- **THEN** the UI shows a local preview without requiring a prior OSS upload

### Requirement: Either side of a gallery pair may be empty
In the admin portfolio pair gallery editor, each pair SHALL allow an empty render (effect) image and an empty real (site) image. A pair with at least one non-empty side SHALL be valid to save. Cover image requirements are unchanged by this capability and SHALL remain required as today.

#### Scenario: Save pair with only real image
- **WHEN** an operator fills only the real slot of a pair and leaves render empty
- **THEN** the pair is accepted on save
- **AND** the submitted `reals` entry at that index is non-empty
- **AND** the submitted `renders` entry at that index is an empty dual placeholder

#### Scenario: Save pair with only render image
- **WHEN** an operator fills only the render slot and leaves real empty
- **THEN** the pair is accepted on save with an empty real placeholder at the same index

#### Scenario: Cover still required
- **WHEN** an operator saves a portfolio
- **THEN** cover required behavior is unchanged by this gallery-pair rule change

### Requirement: Silently drop pairs with both sides empty
On save, pairs where both render and real are empty SHALL be omitted from the submitted `renders` and `reals` arrays without blocking the save or showing a validation error for that empty pair.

#### Scenario: Empty pair skipped
- **WHEN** the editor contains a pair with both sides empty among other valid pairs
- **THEN** save succeeds
- **AND** that empty pair does not appear as an index in the submitted arrays

### Requirement: Persist equal-length aligned arrays
Valid pairs SHALL still be submitted as equal-length `renders` and `reals` arrays in pair order, using empty dual placeholders for the missing side.

#### Scenario: Mixed optional sides keep alignment
- **WHEN** the operator saves pairs that are render-only, real-only, and both-filled in that order
- **THEN** `renders.length` equals `reals.length`
- **AND** index 0 has empty real placeholder, index 1 has empty render placeholder, index 2 has both non-empty

### Requirement: Editor copy reflects optional sides
The pair editor UI SHALL label both render and real slots as optional (or otherwise clearly indicate neither side is mandatory within a pair).

#### Scenario: Labels show optional
- **WHEN** an operator views a gallery pair row
- **THEN** neither slot is presented as the sole required field of the pair
