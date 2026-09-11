export const PORTFOLIO_CATEGORIES = [
  { value: 'residential', label: '住宅', group: '室内设计' },
  { value: 'commercial', label: '商业', group: '室内设计' },
  { value: 'office', label: '办公', group: '室内设计' },
  { value: 'architecture', label: '建筑设计' },
  { value: 'installation', label: '装置设计' },
]

export const PORTFOLIO_CATEGORY_VALUES = PORTFOLIO_CATEGORIES.map((c) => c.value)

export function categoryLabel(value) {
  return PORTFOLIO_CATEGORIES.find((c) => c.value === value)?.label || value
}

export function isValidCategory(value) {
  return PORTFOLIO_CATEGORY_VALUES.includes(value)
}
