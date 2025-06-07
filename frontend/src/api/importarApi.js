import axios from 'axios';

const baseUrl = 'http://localhost:8080';

export const listImports = async (page = 1, limit = 10, token) => {
  try {
    const response = await axios.get(`${baseUrl}/api/imports?page=${page}&limit=${limit}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const uploadExcelFile = async (file, token, onUploadProgress) => {
  try {
    const formData = new FormData();
    formData.append('file', file);

    const response = await axios.post(`${baseUrl}/api/import`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`,
      },
      onUploadProgress,
    });
    return response.data;
  } catch (error) {
    throw error;
  }
}; 

export const getBillingSummaryByResources = async (token) => {
  try {
    const response = await axios.get(`${baseUrl}/api/billing/summary/resources`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getBillingSummaryByClients = async (token) => {
  try {
    const response = await axios.get(`${baseUrl}/api/billing/summary/clients`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getBillingSummaryByMonths = async (token) => {
  try {
    const response = await axios.get(`${baseUrl}/api/billing/summary/months`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
}; 