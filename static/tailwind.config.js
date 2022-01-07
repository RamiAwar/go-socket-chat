const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["../html/*.jet"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
