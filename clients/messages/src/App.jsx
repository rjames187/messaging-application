import Login from "./pages/Login"
import Signup from "./pages/Signup"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/view" element={<View />} />
        <Route path="/update" element={<Update />} />
      </Routes>
    </Router>
  )
}

export default App
