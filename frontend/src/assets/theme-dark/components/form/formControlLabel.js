/**
=========================================================
* Material Dashboard 2 React - v2.2.0
=========================================================

* Product Page:https://erasmocs.netlify.app/
* Copyright 2025 Erasmo Cardoso (https://erasmocs.netlify.app/)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

// Material Dashboard 2 React base styles
import colors from "assets/theme-dark/base/colors";
import typography from "assets/theme-dark/base/typography";

// Material Dashboard 2 React helper functions
import pxToRem from "assets/theme-dark/functions/pxToRem";

const { white } = colors;
const { size, fontWeightBold } = typography;

const formControlLabel = {
  styleOverrides: {
    root: {
      display: "block",
      minHeight: pxToRem(24),
      marginBottom: pxToRem(2),
    },

    label: {
      display: "inline-block",
      fontSize: size.sm,
      fontWeight: fontWeightBold,
      color: white.main,
      lineHeight: 1,
      transform: `translateY(${pxToRem(1)})`,
      marginLeft: pxToRem(4),

      "&.Mui-disabled": {
        color: white.main,
      },
    },
  },
};

export default formControlLabel;
