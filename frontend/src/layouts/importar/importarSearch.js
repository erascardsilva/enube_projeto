import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { listCategories } from '../../api/categoryApi';
import { listResources } from '../../api/resourceApi';
import { listBillings, getBillingSummaryByCategory } from '../../api/billingApi';
import { getBillingSummaryByResources } from '../../api/importarApi';
import { getBillingSummaryByClients } from '../../api/importarApi';
import { getBillingSummaryByMonths } from '../../api/importarApi';

import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Alert from "@mui/material/Alert";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDInput from "components/MDInput";
import MDTypography from "components/MDTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MDPagination from "components/MDPagination";
import DataTable from "examples/Tables/DataTable";

const ImportarSearch = () => {
  // State for Categories
  const [categories, setCategories] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [alert, setAlert] = useState(null);
  const [hasSearched, setHasSearched] = useState(false);

  // State for Resources
  const [resources, setResources] = useState([]);
  const [currentPageResources, setCurrentPageResources] = useState(1);
  const [totalPagesResources, setTotalPagesResources] = useState(1);
  const [hasSearchedResources, setHasSearchedResources] = useState(false);

  // State for Billings
  const [billings, setBillings] = useState([]);
  const [currentPageBillings, setCurrentPageBillings] = useState(1);
  const [totalPagesBillings, setTotalPagesBillings] = useState(1);
  const [hasSearchedBillings, setHasSearchedBillings] = useState(false);

  // State for Billing Summary by Category
  const [billingSummaryCategories, setBillingSummaryCategories] = useState([]);
  const [hasSearchedBillingSummaryCategories, setHasSearchedBillingSummaryCategories] = useState(false);

  // State for Billing Summary by Resources
  const [billingSummaryResources, setBillingSummaryResources] = useState([]);
  const [hasSearchedBillingSummaryResources, setHasSearchedBillingSummaryResources] = useState(false);

  // State for Billing Summary by Clients
  const [billingSummaryClients, setBillingSummaryClients] = useState([]);
  const [hasSearchedBillingSummaryClients, setHasSearchedBillingSummaryClients] = useState(false);

  // State for Billing Summary by Months
  const [billingSummaryMonths, setBillingSummaryMonths] = useState([]);
  const [hasSearchedBillingSummaryMonths, setHasSearchedBillingSummaryMonths] = useState(false);

  const [rowsPerPage] = useState(10);
  const navigate = useNavigate();

  const fetchCategories = async (page, limit) => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await listCategories(page, limit, token);
      if (!data.data || data.data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setCategories([]);
        return;
      }
      setCategories(data.data);
      setTotalPages(data.pagination.total_pages);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setCategories([]);
    }
  };

  const fetchResources = async (page, limit) => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await listResources(page, limit, token);
      if (!data.data || data.data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setResources([]);
        return;
      }
      setResources(data.data);
      setTotalPagesResources(data.pagination.total_pages);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setResources([]);
    }
  };

  const fetchBillings = async (page, limit) => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await listBillings(page, limit, token);
      if (!data.data || data.data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setBillings([]);
        return;
      }
      setBillings(data.data);
      setTotalPagesBillings(data.pagination.total_pages);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setBillings([]);
    }
  };

  const fetchBillingSummaryCategories = async () => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await getBillingSummaryByCategory(token);
      if (!data || data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setBillingSummaryCategories([]);
        return;
      }
      setBillingSummaryCategories(data);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setBillingSummaryCategories([]);
    }
  };

  const fetchBillingSummaryResources = async () => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await getBillingSummaryByResources(token);
      if (!data || data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setBillingSummaryResources([]);
        return;
      }
      setBillingSummaryResources(data);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setBillingSummaryResources([]);
    }
  };

  const fetchBillingSummaryClients = async () => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await getBillingSummaryByClients(token);
      if (!data || data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setBillingSummaryClients([]);
        return;
      }
      setBillingSummaryClients(data);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setBillingSummaryClients([]);
    }
  };

  const fetchBillingSummaryMonths = async () => {
    try {
      const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
      if (!token) {
        setAlert({
          message: "No authentication token found. Please log in.",
          color: "error",
        });
        return;
      }
      const data = await getBillingSummaryByMonths(token);
      if (!data || data.length === 0) {
        setAlert({
          message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
          color: "warning",
        });
        setBillingSummaryMonths([]);
        return;
      }
      setBillingSummaryMonths(data);
    } catch (error) {
      setAlert({
        message: "Nenhum dado encontrado. Por favor, importe o arquivo primeiro.",
        color: "warning",
      });
      setBillingSummaryMonths([]);
    }
  };

  useEffect(() => {
    if (hasSearched) {
      fetchCategories(currentPage, rowsPerPage);
    }
  }, [currentPage, rowsPerPage, hasSearched]);

  useEffect(() => {
    if (hasSearchedResources) {
      fetchResources(currentPageResources, rowsPerPage);
    }
  }, [currentPageResources, rowsPerPage, hasSearchedResources]);

  useEffect(() => {
    if (hasSearchedBillings) {
      fetchBillings(currentPageBillings, rowsPerPage);
    }
  }, [currentPageBillings, rowsPerPage, hasSearchedBillings]);

  useEffect(() => {
    if (hasSearchedBillingSummaryCategories) {
      fetchBillingSummaryCategories();
    }
  }, [hasSearchedBillingSummaryCategories]);

  useEffect(() => {
    if (hasSearchedBillingSummaryResources) {
      fetchBillingSummaryResources();
    }
  }, [hasSearchedBillingSummaryResources]);

  useEffect(() => {
    if (hasSearchedBillingSummaryClients) {
      fetchBillingSummaryClients();
    }
  }, [hasSearchedBillingSummaryClients]);

  useEffect(() => {
    if (hasSearchedBillingSummaryMonths) {
      fetchBillingSummaryMonths();
    }
  }, [hasSearchedBillingSummaryMonths]);

  const handlePageChange = (event, value) => {
    setCurrentPage(value);
  };

  const handleSearchButtonClick = () => {
    setCurrentPage(1); // Reset to first page on new search
    setHasSearched(true);
  };

  const handleResourcePageChange = (event, value) => {
    setCurrentPageResources(value);
  };

  const handleResourceSearchButtonClick = () => {
    setCurrentPageResources(1);
    setHasSearchedResources(true);
  };

  const handleBillingPageChange = (event, value) => {
    setCurrentPageBillings(value);
  };

  const handleBillingSearchButtonClick = () => {
    setCurrentPageBillings(1);
    setHasSearchedBillings(true);
  };

  const handleBillingSummaryCategorySearchButtonClick = () => {
    setHasSearchedBillingSummaryCategories(true);
  };

  const handleBillingSummaryResourceSearchButtonClick = () => {
    setHasSearchedBillingSummaryResources(true);
  };

  const handleBillingSummaryClientSearchButtonClick = () => {
    setHasSearchedBillingSummaryClients(true);
  };

  const handleBillingSummaryMonthSearchButtonClick = () => {
    setHasSearchedBillingSummaryMonths(true);
  };

  const categoryColumns = [
    { Header: "ID", accessor: "id", width: "5%" },
    { Header: "Name", accessor: "name", width: "25%" },
    { Header: "Sub Category", accessor: "sub_category", width: "30%" },
    { Header: "Type", accessor: "type", width: "40%" },
  ];

  const resourceColumns = [
    { Header: "ID", accessor: "ID", width: "5%" },
    { Header: "Name", accessor: "Name", width: "35%" },
    { Header: "Location", accessor: "Location", width: "20%" },
    { Header: "Group", accessor: "Group", width: "20%" },
    { Header: "Consumed Service", accessor: "ConsumedService", width: "20%" },
  ];

  const billingColumns = [
    { Header: "ID", accessor: "id", width: "5%" },
    { Header: "Amount", accessor: "amount", width: "10%" },
    { Header: "Billing Date", accessor: "billing_date", width: "15%" },
    { Header: "Category", accessor: "category_name", width: "15%" },
    { Header: "Client", accessor: "client_name", width: "15%" },
    { Header: "Description", accessor: "description", width: "20%" },
    { Header: "Resource", accessor: "resource_name", width: "20%" },
    { Header: "Unit Price", accessor: "unit_price", width: "10%" },
  ];

  const billingSummaryCategoryColumns = [
    { Header: "Category", accessor: "category", width: "40%" },
    { Header: "Total", accessor: "total", width: "30%" },
    { Header: "Count", accessor: "count", width: "30%" },
  ];

  const billingSummaryResourceColumns = [
    { Header: "Resource Name", accessor: "resource_name", width: "60%" },
    { Header: "Total", accessor: "total", width: "20%" },
    { Header: "Count", accessor: "count", width: "20%" },
  ];

  const billingSummaryClientColumns = [
    { Header: "Client Name", accessor: "client_name", width: "60%" },
    { Header: "Total", accessor: "total", width: "20%" },
    { Header: "Count", accessor: "count", width: "20%" },
  ];

  const billingSummaryMonthColumns = [
    { Header: "Month", accessor: "month", width: "30%" },
    { Header: "Total", accessor: "total", width: "30%" },
    { Header: "Count", accessor: "count", width: "40%" },
  ];

  return (
    <DashboardLayout>
      <DashboardNavbar showTitle={false} />
      <MDBox pt={6} pb={3}>
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
                  Listar todas as categorias
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearched && categories.length > 0 ? (
                  <DataTable
                    table={{ columns: categoryColumns, rows: categories }}
                  />
                ) : hasSearched && categories.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhuma categoria encontrada.
                  </MDTypography>
                ) : null}

                {hasSearched && totalPages > 1 && (
                  <MDBox mt={3} display="flex" justifyContent="center">
                    <MDPagination
                      count={totalPages}
                      page={currentPage}
                      onChange={handlePageChange}
                    />
                  </MDBox>
                )}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Listar todos os recursos (com paginação)
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleResourceSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedResources && resources.length > 0 ? (
                  <DataTable
                    table={{ columns: resourceColumns, rows: resources }}
                  />
                ) : hasSearchedResources && resources.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum recurso encontrado.
                  </MDTypography>
                ) : null}

                {hasSearchedResources && totalPagesResources > 1 && (
                  <MDBox mt={3} display="flex" justifyContent="center">
                    <MDPagination
                      count={totalPagesResources}
                      page={currentPageResources}
                      onChange={handleResourcePageChange}
                    />
                  </MDBox>
                )}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Listar todos os faturamentos (com paginação)
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleBillingSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedBillings && billings.length > 0 ? (
                  <DataTable
                    table={{ columns: billingColumns, rows: billings }}
                  />
                ) : hasSearchedBillings && billings.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum faturamento encontrado.
                  </MDTypography>
                ) : null}

                {hasSearchedBillings && totalPagesBillings > 1 && (
                  <MDBox mt={3} display="flex" justifyContent="center">
                    <MDPagination
                      count={totalPagesBillings}
                      page={currentPageBillings}
                      onChange={handleBillingPageChange}
                    />
                  </MDBox>
                )}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Resumo de faturamento por categoria
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleBillingSummaryCategorySearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedBillingSummaryCategories && billingSummaryCategories.length > 0 ? (
                  <DataTable
                    table={{ columns: billingSummaryCategoryColumns, rows: billingSummaryCategories }}
                  />
                ) : hasSearchedBillingSummaryCategories && billingSummaryCategories.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum resumo de faturamento por categoria encontrado.
                  </MDTypography>
                ) : null}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Resumo de faturamento por recursos
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleBillingSummaryResourceSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedBillingSummaryResources && billingSummaryResources.length > 0 ? (
                  <DataTable
                    table={{ columns: billingSummaryResourceColumns, rows: billingSummaryResources }}
                  />
                ) : hasSearchedBillingSummaryResources && billingSummaryResources.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum resumo de faturamento por recursos encontrado.
                  </MDTypography>
                ) : null}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Resumo de faturamento por clientes
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleBillingSummaryClientSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedBillingSummaryClients && billingSummaryClients.length > 0 ? (
                  <DataTable
                    table={{ columns: billingSummaryClientColumns, rows: billingSummaryClients }}
                  />
                ) : hasSearchedBillingSummaryClients && billingSummaryClients.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum resumo de faturamento por clientes encontrado.
                  </MDTypography>
                ) : null}
              </MDBox>
            </Card>
          </Grid>
          <Grid item xs={12} mt={3}>
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
                  Resumo de faturamento por meses
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
                  <MDButton
                    variant="gradient"
                    color="info"
                    onClick={handleBillingSummaryMonthSearchButtonClick}
                  >
                    Buscar
                  </MDButton>
                </MDBox>
                {hasSearchedBillingSummaryMonths && billingSummaryMonths.length > 0 ? (
                  <DataTable
                    table={{ columns: billingSummaryMonthColumns, rows: billingSummaryMonths }}
                  />
                ) : hasSearchedBillingSummaryMonths && billingSummaryMonths.length === 0 ? (
                  <MDTypography variant="body2" color="text.secondary" mt={2} ml={2}>
                    Nenhum resumo de faturamento por meses encontrado.
                  </MDTypography>
                ) : null}
              </MDBox>
            </Card>
          </Grid>
        </Grid>
      </MDBox>
      <Footer />
    </DashboardLayout>
  );
};

export default ImportarSearch; 