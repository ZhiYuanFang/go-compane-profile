export const PORTFOLIO_CATEGORIES = [
  { value: 'residential', label: '住宅' },
  { value: 'commercial', label: '商业' },
  { value: 'office', label: '办公' },
  { value: 'installation', label: '装置' },
]

export const PORTFOLIO_CATEGORY_VALUES = PORTFOLIO_CATEGORIES.map((c) => c.value)

export function categoryLabel(value) {
  return PORTFOLIO_CATEGORIES.find((c) => c.value === value)?.label || value
}

export function isValidCategory(value) {
  return PORTFOLIO_CATEGORY_VALUES.includes(value)
}
