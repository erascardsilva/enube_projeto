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

// prop-types is a library for typechecking of props
import PropTypes from "prop-types";

// @mui material components
import Card from "@mui/material/Card";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";

// Material Dashboard 2 React components
import { useMaterialUIController } from "context";

// Timeline context
import { TimelineProvider } from "examples/Timeline/context";

function TimelineList({ title, dark, children }) {
  const [controller] = useMaterialUIController();
  const { darkMode } = controller;

  return (
    <TimelineProvider value={dark}>
      <Card>
        <MDBox
          bgColor={dark ? "dark" : "white"}
          variant="gradient"
          borderRadius="xl"
          sx={{ background: ({ palette: { background } }) => darkMode && background.card }}
        >
          <MDBox pt={3} px={3}>
            <MDTypography variant="h6" fontWeight="medium" color={dark ? "white" : "dark"}>
              {title}
            </MDTypography>
          </MDBox>
          <MDBox p={2}>{children}</MDBox>
        </MDBox>
      </Card>
    </TimelineProvider>
  );
}

// Setting default values for the props of TimelineList
TimelineList.defaultProps = {
  dark: false,
};

// Typechecking props for the TimelineList
TimelineList.propTypes = {
  title: PropTypes.string.isRequired,
  dark: PropTypes.bool,
  children: PropTypes.node.isRequired,
};

export default TimelineList;
