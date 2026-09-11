## ADDED Requirements

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
