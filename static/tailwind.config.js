const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["../html/*.jet"],
  theme: {
    extend: {
      borderRadius: {
        "4xl": "1.9rem"
      },
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
