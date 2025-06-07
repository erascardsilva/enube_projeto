import React, { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { uploadExcelFile } from '../../api/importarApi';

import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Alert from "@mui/material/Alert";
import LinearProgress from "@mui/material/LinearProgress";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";

import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDInput from "components/MDInput";
import MDTypography from "components/MDTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";

const ImportarAdd = () => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [alert, setAlert] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const [isConverting, setIsConverting] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [conversionMessage, setConversionMessage] = useState("");
  const [responseDetails, setResponseDetails] = useState(null);
  const [elapsedTime, setElapsedTime] = useState(0);
  const timerIntervalRef = useRef(null);
  const navigate = useNavigate();

  const handleFileChange = (event) => {
    setSelectedFile(event.target.files[0]);
    setAlert(null);
    setIsUploading(false);
    setIsConverting(false);
    setUploadProgress(0);
    setConversionMessage("");
    setResponseDetails(null);
    setElapsedTime(0);
    if (timerIntervalRef.current) {
      clearInterval(timerIntervalRef.current);
      timerIntervalRef.current = null;
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setAlert({ message: "Por favor, selecione um arquivo Excel para importar.", color: "error" });
      return;
    }

    setIsUploading(true);
    setAlert(null);
    setUploadProgress(0);
    setConversionMessage("Iniciando upload...");
    setResponseDetails(null);
    setElapsedTime(0);
    timerIntervalRef.current = setInterval(() => {
      setElapsedTime(prevTime => prevTime + 1);
    }, 1000);

    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "Token de autenticação não encontrado. Por favor, faça login.",
          color: "error",
        });
        setIsUploading(false);
        setConversionMessage("");
        if (timerIntervalRef.current) {
          clearInterval(timerIntervalRef.current);
          timerIntervalRef.current = null;
        }
        return;
      }

      const data = await uploadExcelFile(selectedFile, token, (progressEvent) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        setUploadProgress(percentCompleted);
        if (percentCompleted === 100) {
          setIsUploading(false);
          setIsConverting(true);
          setConversionMessage("Upload concluído. Processando e convertendo dados...");
        }
      });

      setAlert({ message: data.message, color: "success" });
      setResponseDetails(data.stats);
      setSelectedFile(null);
    } catch (error) {
      console.error('Erro ao importar arquivo:', error);
      setAlert({
        message: "Erro ao importar arquivo: " + (error.response?.data?.message || error.message),
        color: "error",
      });
    } finally {
      setIsUploading(false);
      setIsConverting(false);
      setConversionMessage("");
      setUploadProgress(0);
      if (timerIntervalRef.current) {
        clearInterval(timerIntervalRef.current);
        timerIntervalRef.current = null;
      }
    }
  };

  const formatTime = (totalSeconds) => {
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <DashboardLayout>
      <DashboardNavbar showTitle={false} />
      <MDBox pt={6} pb={3}>
        <Grid container spacing={3} justifyContent="center">
          <Grid item xs={12} lg={8}>
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
                  Importar Arquivo Excel do teste...
                </MDTypography>
              </MDBox>
              <MDBox p={3}>
                {alert && (
                  <Alert
                    severity={alert.color === "error" ? "error" : "success"}
                    onClose={() => setAlert(null)}
                  >
                    {alert.message}
                  </Alert>
                )}
                <MDBox display="flex" alignItems="center" gap={2} mb={2}>
                  <MDInput
                    type="file"
                    inputProps={{ accept: ".xlsx, .xls" }}
                    onChange={handleFileChange}
                    fullWidth
                  />
                </MDBox>
                <MDBox display="flex" alignItems="center" gap={2}>
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleUpload}
                    disabled={isUploading || isConverting}
                  >
                    {(isUploading || isConverting) ? "Aguarde..." : "Upload Arquivo"}
                  </MDButton>
                  {(isUploading || isConverting || responseDetails) && (
                    <MDTypography variant="body2" color="text.secondary">
                      Tempo Decorrido: {formatTime(elapsedTime)}
                    </MDTypography>
                  )}
                </MDBox>

                {(isUploading || isConverting) && (
                  <MDBox mt={2} mb={2}>
                    {isUploading && (
                      <MDBox sx={{ overflow: 'hidden' }}>
                        <MDTypography variant="body2" color="text.secondary" mb={1}>
                          Progresso do Upload: {uploadProgress}%
                        </MDTypography>
                        <LinearProgress variant="determinate" value={uploadProgress} color="info" />
                      </MDBox>
                    )}
                    {isConverting && (
                      <MDBox mt={isUploading ? 2 : 0} sx={{ overflow: 'hidden' }}>
                        <MDTypography variant="body2" color="text.secondary" mb={1}>
                          Progresso da Conversão:
                        </MDTypography>
                        <LinearProgress variant="indeterminate" color="info" />
                      </MDBox>
                    )}
                    {conversionMessage && (
                      <MDTypography variant="body2" color="text.secondary" mt={1}>
                        {conversionMessage}
                      </MDTypography>
                    )}
                  </MDBox>
                )}

                {responseDetails && (
                  <MDBox mt={4}>
                    <MDTypography variant="h6" mb={2}>Detalhes da Importação:</MDTypography>
                    <TableContainer component={Paper}>
                      <Table size="small">
                        <TableBody>
                          <TableRow>
                            <TableCell>Header Rows</TableCell>
                            <TableCell>{responseDetails.header_rows}</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>Imported Rows</TableCell>
                            <TableCell>{responseDetails.imported_rows}</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>Skipped Rows</TableCell>
                            <TableCell>{responseDetails.skipped_rows}</TableCell>
                          </TableRow>
                          <TableRow>
                            <TableCell>Total Rows</TableCell>
                            <TableCell>{responseDetails.total_rows}</TableCell>
                          </TableRow>
                        </TableBody>
                      </Table>
                    </TableContainer>

                    <MDTypography variant="h6" mt={3} mb={1}>Categorias:</MDTypography>
                    <TableContainer component={Paper}>
                      <Table size="small">
                        <TableHead>
                          <TableRow>
                            <TableCell>Categoria</TableCell>
                            <TableCell>Contagem</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {Object.entries(responseDetails.categories).map(([category, count]) => (
                            <TableRow key={category}>
                              <TableCell>{category}</TableCell>
                              <TableCell>{count}</TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </MDBox>
                )}
              </MDBox>
            </Card>
          </Grid>
        </Grid>
      </MDBox>
      <Footer />
    </DashboardLayout>
  );
};

export default ImportarAdd; 