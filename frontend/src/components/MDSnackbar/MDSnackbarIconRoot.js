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

// @mui material components
import Icon from "@mui/material/Icon";
import { styled } from "@mui/material/styles";

export default styled(Icon)(({ theme, ownerState }) => {
  const { palette, functions, typography } = theme;
  const { color, bgWhite } = ownerState;

  const { white, transparent, gradients } = palette;
  const { pxToRem, linearGradient } = functions;
  const { size } = typography;

  // backgroundImage value
  let backgroundImageValue;

  if (bgWhite) {
    backgroundImageValue = gradients[color]
      ? linearGradient(gradients[color].main, gradients[color].state)
      : linearGradient(gradients.info.main, gradients.info.state);
  } else if (color === "light") {
    backgroundImageValue = linearGradient(gradients.dark.main, gradients.dark.state);
  }

  return {
    backgroundImage: backgroundImageValue,
    WebkitTextFillColor: bgWhite || color === "light" ? transparent.main : white.main,
    WebkitBackgroundClip: "text",
    marginRight: pxToRem(8),
    fontSize: size.lg,
    transform: `translateY(${pxToRem(-2)})`,
  };
});
