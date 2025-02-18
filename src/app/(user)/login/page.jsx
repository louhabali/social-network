'use client';
import '../../styles/login.css'
import { useState } from "react"
import Link from "next/link"
export default function loginPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }
  const handleSubmit = async (e) => {
    e.preventDefault()
    console.log(formData)
  
    const response = await fetch("/api/login", {
      method: "POST",
      body: JSON.stringify(formData),
    })

    const result = await response.json()
    console.log(result)
  }
  return (
    <div className="container">
      <form onSubmit={handleSubmit} className="form">
        <input type="email" name="email" placeholder="Email" required className="inputField" onChange={handleChange} />
        <input type="password" name="password" placeholder="Password" className="inputField" required onChange={handleChange} />
        <button type="submit" className="submitButton">login</button>
      </form>
      <p className="link">
        Create an account? <Link href="/register" className='registerlink'>Register</Link>
      </p>
    </div>
  )
}
