/**
 * Map items with bounded concurrency, preserving input order in the result.
 * @template T, R
 * @param {T[]} items
 * @param {number} concurrency
 * @param {(item: T, index: number) => Promise<R>} fn
 * @returns {Promise<R[]>}
 */
export async function mapPool(items, concurrency, fn) {
  const list = Array.isArray(items) ? items : []
  const n = list.length
  if (!n) return []

  const limit = Math.max(1, Math.min(concurrency || 1, n))
  const results = new Array(n)
  let nextIndex = 0

  async function worker() {
    while (true) {
      const i = nextIndex++
      if (i >= n) return
      results[i] = await fn(list[i], i)
    }
  }

  const workers = []
  for (let w = 0; w < limit; w++) workers.push(worker())
  await Promise.all(workers)
  return results
}
