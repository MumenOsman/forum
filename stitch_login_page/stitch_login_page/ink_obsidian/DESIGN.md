# Design System: The Illuminated Manuscript

## 1. Overview & Creative North Star
**Creative North Star: "The Modern Archive"**
This design system moves away from the sterile, utilitarian "forum" look of the early 2000s and instead draws inspiration from high-end editorial magazines and digital archives. It is a space where the weight of classic literature meets the fluidity of modern technology.

We reject the "boxed" layout. Instead of rigid grids and heavy borders, we utilize **Intentional Asymmetry** and **Tonal Depth**. By overlapping glass-morphic elements and using high-contrast typography scales, we create a sense of physical layering—as if digital vellum were resting atop a deep, ink-stained desk. The goal is to make the user feel like they are contributing to a lasting collection of thought, rather than a fleeting chat thread.

---

## 2. Colors & Surface Philosophy
The palette is rooted in the depth of a midnight library (`background: #0b1326`), punctuated by the "electric gold" of a scholar’s lamp (`tertiary: #fabd00`).

### The "No-Line" Rule
Traditional 1px solid borders are strictly prohibited for sectioning. Definition must be achieved through:
*   **Background Shifts:** Use `surface_container_low` against `surface` to define a sidebar.
*   **Tonal Transitions:** Allow the `surface_container` hierarchy to create natural boundaries.
*   **Negative Space:** Rely on the Spacing Scale (specifically `8` and `12`) to create clear "islands" of content.

### Surface Hierarchy & Nesting
Treat the UI as a series of nested physical layers. 
*   **Base:** `surface_dim` (#0b1326)
*   **Lower Content:** `surface_container_low` (#131b2e)
*   **Standard Cards:** `surface_container` (#171f33)
*   **Promoted/Active Cards:** `surface_container_high` (#222a3d)

### The "Glass & Gradient" Rule
To evoke a premium feel, floating navigation and hero elements must utilize **Glassmorphism**:
*   **Background:** `rgba(23, 31, 51, 0.7)` (Surface Container)
*   **Effect:** `backdrop-filter: blur(12px) saturate(180%);`
*   **Signature Gradient:** For primary CTAs, transition from `primary` (#9ecaff) to `primary_container` (#003157) at a 135-degree angle. This adds a "soul" to the interactive elements that flat colors lack.

---

### 3. Typography
We use a high-contrast pairing to balance "The Literary" with "The Functional."

*   **The Voice (Serif - Newsreader):** Used for all `display` and `headline` levels. It conveys authority and history. The large x-height of Newsreader ensures legibility even at high weights.
*   **The Engine (Sans-Serif - Manrope):** Used for `title`, `body`, and `label`. Manrope is a modern geometric sans that provides a clean, technical counter-balance to the serif headers, ensuring long forum posts remain highly readable.

**Hierarchy Goal:** A `display-lg` headline should feel like a book title, while `body-md` feels like a well-set page of prose.

---

## 4. Elevation & Depth
Depth is not an effect; it is information.

*   **The Layering Principle:** Avoid shadows for static elements. Instead, place a `surface_container_lowest` (#060e20) card inside a `surface_container_high` (#222a3d) parent to create a "recessed" look.
*   **Ambient Shadows:** For floating modals or dropdowns, use a "scholar's shadow": `box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);`. The shadow color must never be pure black; it should be a deep tint of the surface color to maintain a natural, atmospheric feel.
*   **The "Ghost Border" Fallback:** When a container sits on a background of similar value, use a `1px` border of `outline_variant` at **15% opacity**. This provides a "whisper" of an edge without breaking the glass aesthetic.

---

## 5. Components

### Cards & Discussion Threads
*   **Forbid Dividers:** Do not use `<hr>` or border-bottoms between posts. Use a background shift (e.g., alternating `surface_container` and `surface_container_low`) or a `spacing-6` gap.
*   **Layout:** Use asymmetrical padding—more generous on the "reading side" (left) to allow the text to breathe.

### Buttons
*   **Primary:** `surface_tint` background with `on_primary_container` text. `rounded-md` (0.375rem). Use a subtle inner-glow (1px top border, 20% white) to simulate a glass edge.
*   **Tertiary (Action):** `tertiary` text on no background. Use `label-md` for uppercase styling to denote a "meta" action (e.g., "REPLY", "SHARE").

### Inputs & Fields
*   **Style:** Minimalist. No background fill, only a `ghost border` on the bottom. On focus, the border transitions to `primary` and a subtle `surface_container_highest` background fades in.
*   **Typography:** All user-input text should be `body-lg` in Manrope.

### Chips (Tags)
*   **Category Chips:** Use `secondary_container` with `on_secondary_container` text. Use `rounded-full` for a soft, pill-shaped organic feel.

---

## 6. Do's and Don'ts

### Do:
*   **Embrace the Blur:** Use `backdrop-filter` on any element that overlays another (modals, sticky navs).
*   **Mix Weights:** Pair a `display-sm` (Serif) with a `label-sm` (Sans-Serif) metadata tag immediately above it.
*   **Use Intentional White Space:** If a section feels "messy," add more padding rather than adding a border.

### Don't:
*   **Don't use pure white:** The brightest text should be `on_surface` (#dae2fd). Pure white (#ffffff) is too harsh for this "ink and vellum" aesthetic.
*   **Don't use standard "Drop Shadows":** Avoid the "fuzzy grey" look. Shadows should be wide, soft, and deep.
*   **Don't use Framework Defaults:** Avoid the 8px rounded corner typical of Material or Bootstrap. Stick to our defined `rounded-md` (0.375rem) or `rounded-xl` (0.75rem) for a more bespoke architectural feel.