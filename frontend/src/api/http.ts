import axios from 'axios'

// 创建 axios 实例，后续所有 API 都从这里走。
// 这样方便统一处理 baseURL、超时、拦截器。
const http = axios.create({
  baseURL: 'http://192.168.64.129:8080/api',
  timeout: 10000,
})

// 响应拦截器：统一返回 data，减少页面层处理负担。
http.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 这里尽量把后端报错信息透出来，方便你调试
    const msg =
      error?.response?.data?.error ||
      error?.response?.data?.message ||
      error?.message ||
      'request failed'
    return Promise.reject(new Error(msg))
  },
)

export default http