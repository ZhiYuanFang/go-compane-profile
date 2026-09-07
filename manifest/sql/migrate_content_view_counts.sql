-- Add view counters for portfolio, activity, and about-us.
ALTER TABLE `portfolio`
  ADD COLUMN `view_count` INT NOT NULL DEFAULT 0 AFTER `cover_thumb_url`;

ALTER TABLE `activity`
  ADD COLUMN `view_count` INT NOT NULL DEFAULT 0 AFTER `sort_order`;

ALTER TABLE `company`
  ADD COLUMN `about_view_count` INT NOT NULL DEFAULT 0 AFTER `logo_hor_thumb_url`;
