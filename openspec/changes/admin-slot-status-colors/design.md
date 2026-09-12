## Context

Slot statuses already map to Chinese labels via `dualStatusLabel`. Gallery badges and cover status use almost one accent color (plus error red). Product locked **four tiers** and wants **border/outline** as well as text color.

## Goals / Non-Goals

**Goals:**

- Four tones: `wait` | `busy` | `done` | `error`
- Status text and cell/field border use the same tone
- Gallery grid cells + DualImageField cover share the mapping

**Non-Goals:**

- Changing status copy or upload pipeline
- Icons/animations beyond color
- Light-theme palette redesign

## Decisions

1. **Tone mapping**  
   - `queued` → `wait` (灰蓝)  
   - `compressing` | `uploading` → `busy` (accent 蓝或琥珀，选 **accent 蓝** 与现有 busy 一致)  
   - `done` (or remote URLs without pending work) → `done` (绿)  
   - `error` → `error` (现有 danger)  
   - idle / empty → no tone class (default border)

2. **Shared helper**  
   - `dualStatusTone(d)` returns `'' | 'wait' | 'busy' | 'done' | 'error'`  
   - CSS classes: `tone-wait`, `tone-busy`, `tone-done`, `tone-error` on badge and on `.gallery-cell` / `.dual`

3. **Border**  
   - Gallery: cell `border-color` (and slight tint if needed) by tone  
   - Cover: `.dual` border-color by tone when status present  
   - Selected gallery cells: keep selected accent; tone border can be secondary or override when not selected—prefer **tone border always when tone set**, selected may add outline/background separately without fighting tone color

4. **Colors (dark admin)**  
   - wait: `#8b9bb4`  
   - busy: `var(--accent)` / `#6ea8fe`  
   - done: `#3dd68c`  
   - error: `#ff8f7a` / `var(--danger)`

## Risks / Trade-offs

- [Selected + tone both accent] → Selected uses background/check; border stays tone  
- [Color-blind] → Text still present; color is secondary cue  

## Migration Plan

1. Admin JS/CSS only  
2. Sync `resource/public/admin`  
3. Hard-refresh  

## Open Questions

- None — four tiers + border locked.
