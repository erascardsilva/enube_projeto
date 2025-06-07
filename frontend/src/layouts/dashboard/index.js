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
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Icon from "@mui/material/Icon";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";

// Material Dashboard 2 React example components
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";

function Dashboard() {
  const appInfo = [
    { label: "Importados de arquivo", value: "Excel" },
    { label: "Backend", value: "Golang (Gin, GORM, CORS)" },
    { label: "Frontend", value: "React" },
    { label: "Banco de dados", value: "PostgreSQL" },
    { label: "Orquestração", value: "Docker Compose" },
  ];

  const securityInfo = [
    { label: "Autenticação JWT", value: true },
    { label: "Criptografia Bcrypt", value: true },
  ];

  const projectTimeline = [
    { label: "Início do Projeto", value: "05/Junho" },
    { label: "Término do Projeto", value: "11/Junho" },
  ];

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <MDBox py={3}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <MDBox
                mx={2}
                mt={-3}
                py={3}
                px={2}
                variant="gradient"
                bgColor="info"
                borderRadius="lg"
                coloredShadow="info"
              >
                <MDTypography variant="h6" color="white">
                  Informativos do Aplicativo
                </MDTypography>
              </MDBox>
              <MDBox p={3}>
                <MDBox mb={2}>
                  <MDTypography variant="h6" gutterBottom>
                    Visão Geral:
                  </MDTypography>
                  {appInfo.map((item, index) => (
                    <MDBox key={index} display="flex" alignItems="center" mb={1}>
                      <MDTypography variant="body2" color="text.secondary" mr={1}>
                        {item.label}:
                      </MDTypography>
                      <MDTypography variant="body2" color="text.primary">
                        {item.value}
                      </MDTypography>
                    </MDBox>
                  ))}
                </MDBox>

                <MDBox mb={2}>
                  <MDTypography variant="h6" gutterBottom>
                    Segurança:
                  </MDTypography>
                  {securityInfo.map((item, index) => (
                    <MDBox key={index} display="flex" alignItems="center" mb={1}>
                      <MDTypography variant="body2" color="text.secondary" mr={1}>
                        {item.label}:
                      </MDTypography>
                      {item.value ? (
                        <Icon color="success" fontSize="small">check</Icon>
                      ) : (
                        <Icon color="error" fontSize="small">close</Icon>
                      )}
                    </MDBox>
                  ))}
                </MDBox>

                <MDBox mb={2}>
                  <MDTypography variant="h6" gutterBottom>
                    Início do Projeto:
                  </MDTypography>
                  {projectTimeline.map((item, index) => (
                    <MDBox key={index} display="flex" alignItems="center" mb={1}>
                      <MDTypography variant="body2" color="text.secondary" mr={1}>
                        {item.label}:
                      </MDTypography>
                      <MDTypography variant="body2" color="text.primary">
                        {item.value}
                      </MDTypography>
                    </MDBox>
                  ))}
                </MDBox>

                <MDBox mt={4} textAlign="right">
                  <MDTypography variant="h6" color="text.primary">
                    Erasmo Cardoso
                  </MDTypography>
                  <MDTypography variant="body2" color="text.secondary">
                    Desenvolvedor full stack
                  </MDTypography>
                </MDBox>
            </MDBox>
            </Card>
          </Grid>
        </Grid>
      </MDBox>
      <Footer />
    </DashboardLayout>
  );
}

export default Dashboard;
