import Login from "./pages/Login"
import Signup from "./pages/Signup"
import View from "./pages/View"
import Update from "./pages/Update"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import './App.css'

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
