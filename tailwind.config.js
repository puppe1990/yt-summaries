/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./web/templates/**/*.html"],
  safelist: [
    "cais-password-wrap",
    "cais-password-toggle",
    "cais-chat-scroll-down",
    "cais-thinking",
    "cais-thinking-dots",
    "cais-select-search",
    "cais-select-search-native",
    "cais-select-search-trigger",
    "cais-select-search-panel",
    "cais-select-search-input",
    "cais-select-search-list",
    "cais-select-search-option",
    "cais-select-search-label",
    "cais-select-search-chevron",
    "is-selected",
    "is-highlighted",
    "is-hidden",
  ],
  theme: {
    extend: {
      colors: {
        /* CSS variables (input.css) so html.light can flip the palette (#258). */
        ink: "rgb(var(--amarra-ink) / <alpha-value>)",
        foam: "rgb(var(--amarra-foam) / <alpha-value>)",
        copper: "rgb(var(--amarra-copper) / <alpha-value>)",
        tide: "rgb(var(--amarra-tide) / <alpha-value>)",
      },
      fontFamily: {
        sans: ["Source Serif 4", "Georgia", "serif"],
        serif: ["Fraunces", "Georgia", "serif"],
        display: ["Fraunces", "Georgia", "serif"],
        mono: ["IBM Plex Mono", "ui-monospace", "SFMono-Regular", "Menlo", "monospace"],
      },
      boxShadow: {
        "2xs": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
        xs: "0 1px 2px 0 rgb(0 0 0 / 0.05)",
      },
    },
  },
  plugins: [],
};
