## 1. DualImageField interaction

- [x] 1.1 Refactor file handling into a shared `applyFile(File)` used by the file input change handler
- [x] 1.2 Add drag-and-drop on the field root (preventDefault, drag highlight, take first image/* file)
- [x] 1.3 When dropping onto a non-empty preview, prompt for replace confirmation; cancel keeps current value
- [x] 1.4 Empty preview click opens the hidden file input; non-empty preview click keeps lightbox enlarge
- [x] 1.5 Keep the choose/re-choose button working; show clear error for non-image drops

## 2. Build

- [x] 2.1 Build admin and sync static assets to `resource/public/admin` if that is the serving path
