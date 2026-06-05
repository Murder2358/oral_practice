import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ConfigProvider, theme } from 'antd'
import Home from './pages/Home'
import Practice from './pages/Practice'
import Summary from './pages/Summary'

function App() {
  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1677ff',
          borderRadius: 8,
        },
      }}
    >
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/practice" element={<Practice />} />
          <Route path="/summary" element={<Summary />} />
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
