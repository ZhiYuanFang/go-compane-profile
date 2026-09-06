USE `compane_profile`;

-- Company copy from former mini-program mock. Logos left empty.
INSERT INTO `company` (
  `id`,
  `intro_title`,
  `intro_body`,
  `design_values`,
  `years_label`,
  `address`,
  `logo_original_url`,
  `logo_thumb_url`,
  `logo_hor_original_url`,
  `logo_hor_thumb_url`
) VALUES (
  1,
  '目后空间设计 · 公司简介说明',
  '目后空间设计 MUHU DESIGN
是由夏克进先生、方远远女士两位设计师于2016年联合创办成立的设计机构具备开拓精神与国际视野的联合设计机构。

专注于高端住宅系列，秉承着「空间设计是以人为本」这一设计理念，通过整合建筑、空间、品牌、灯光、软装、产品等多方面领域资源，将项目更完整的整合与落地',
  '设计价值观：
尊重每个项目的原创性，赋予其独特的归属与灵魂。',
  '2016-2026 MUHUDESIGN',
  '乐清市总部经济2-703',
  '',
  '',
  '',
  ''
) ON DUPLICATE KEY UPDATE
  `intro_title` = VALUES(`intro_title`),
  `intro_body` = VALUES(`intro_body`),
  `design_values` = VALUES(`design_values`),
  `years_label` = VALUES(`years_label`),
  `address` = VALUES(`address`);

-- Singleton pricing row; URLs empty until admin upload.
INSERT INTO `pricing` (`id`, `original_url`, `thumb_url`)
VALUES (1, '', '')
ON DUPLICATE KEY UPDATE `id` = `id`;
