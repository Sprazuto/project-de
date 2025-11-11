import { isToday } from './helpers'

export const avatarText = (value) => {
  if (!value) return ''
  const nameArray = value.split(' ')

  return nameArray.map((word) => word.charAt(0).toUpperCase()).join('')
}

// TODO: Try to implement this: https://twitter.com/fireship_dev/status/1565424801216311297
export const kFormatter = (num) => {
  const regex = /\B(?=(\d{3})+(?!\d))/g

  return Math.abs(num) > 9999
    ? `${Math.sign(num) * +(Math.abs(num) / 1000).toFixed(1)}k`
    : Math.abs(num).toFixed(0).replace(regex, ',')
}

/**
 * Format and return date in Humanize format
 * Intl docs: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/format
 * Intl Constructor: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl/DateTimeFormat/DateTimeFormat
 * @param {string} value date to format
 * @param {Intl.DateTimeFormatOptions} formatting Intl object to format with
 */
export const formatDate = (value, formatting = { month: 'short', day: 'numeric', year: 'numeric' }) => {
  if (!value) return value

  return new Intl.DateTimeFormat('en-US', formatting).format(new Date(value))
}

/**
 * Return short human friendly month representation of date
 * Can also convert date to only time if date is of today (Better UX)
 * @param {string} value date to format
 * @param {boolean} toTimeForCurrentDay Shall convert to time if day is today/current
 */
export const formatDateToMonthShort = (value, toTimeForCurrentDay = true) => {
  const date = new Date(value)
  let formatting = { month: 'short', day: 'numeric' }
  if (toTimeForCurrentDay && isToday(date)) formatting = { hour: 'numeric', minute: 'numeric' }

  return new Intl.DateTimeFormat('en-US', formatting).format(new Date(value))
}
export const prefixWithPlus = (value) => (value > 0 ? `+${value}` : value)

/**
 * Abbreviate Rupiah values for display (e.g., Rp87.374.471.971,00 -> Rp87.4M)
 * @param {string|number} value Rupiah value as string or number
 * @returns {string} Abbreviated Rupiah string
 */
export const rupiahAbbreviate = (value) => {
  if (!value) return 'Rp0'

  // Parse the Rupiah string to number
  let num
  if (typeof value === 'string') {
    // Remove 'Rp' and replace '.' with '' and ',' with '.' for parsing
    const cleanValue = value.replace('Rp', '').replace(/\./g, '').replace(',', '.')
    num = parseFloat(cleanValue)
  } else {
    num = Number(value)
  }

  if (isNaN(num)) return value

  const absNum = Math.abs(num)
  const sign = num < 0 ? '-' : ''

  if (absNum >= 1e12) {
    return `${sign}Rp${(absNum / 1e12).toFixed(1)}T`
  } else if (absNum >= 1e9) {
    return `${sign}Rp${(absNum / 1e9).toFixed(1)}M`
  } else if (absNum >= 1e6) {
    return `${sign}Rp${(absNum / 1e6).toFixed(1)}Jt`
  } else if (absNum >= 1e3) {
    return `${sign}Rp${(absNum / 1e3).toFixed(1)}K`
  } else {
    return `${sign}Rp${absNum.toLocaleString('id-ID')}`
  }
}
