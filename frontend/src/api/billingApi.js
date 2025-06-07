import axios from 'axios';

const baseUrl = "http://localhost:8080";

export const listBillings = async (page = 1, limit = 10) => {
  try {
    const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
    if (!token) {
      throw new Error("No authentication token found. Please log in.");
    }
    const response = await axios.get(`${baseUrl}/api/billing`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      params: {
        page,
        limit,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching billings:', error);
    throw error;
  }
};

export const getBillingSummaryByCategory = async () => {
  try {
    const token = localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).token : null;
    if (!token) {
      throw new Error("No authentication token found. Please log in.");
    }
    const response = await axios.get(`${baseUrl}/api/billing/summary/categories`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error('Error fetching billing summary by category:', error);
    throw error;
  }
}; 