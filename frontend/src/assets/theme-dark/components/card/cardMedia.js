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

// Material Dashboard 2 React Base Styles
import borders from "assets/theme-dark/base/borders";

// Material Dashboard 2 React Helper Functions
import pxToRem from "assets/theme-dark/functions/pxToRem";

const { borderRadius } = borders;

const cardMedia = {
  styleOverrides: {
    root: {
      borderRadius: borderRadius.xl,
      margin: `${pxToRem(16)} ${pxToRem(16)} 0`,
    },

    media: {
      width: "auto",
    },
  },
};

export default cardMedia;
