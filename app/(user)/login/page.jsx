'use client';
import '../../styles/login.css'
import { useRouter } from "next/navigation"
import { useState } from "react"
import Link from "next/link"
import InvalidCredentials from '@/components/invalidCredentials'

export default function loginPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })
  const [errorMessage, setErrorMessage] = useState("");

  const router = useRouter()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }
  const handleSubmit = async (e) => {
    e.preventDefault()
    if (formData.email.trim().length==0 || formData.password.trim().length==0) {
      setErrorMessage("Please fill all the required fields")
      return
    }
    
    try {
      const response = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        body: JSON.stringify(formData),
        credentials: "include",
      })

      const result = await response.text()

      if (response.ok) {
        router.push("/")
      } else {
        setErrorMessage(result)
      }
    } catch (error) {
      console.error('Error', error)
    }
  }
  return (
    <div className="container">
      <form onSubmit={handleSubmit} className="form">
        <input type="email" name="email" placeholder="Email" required className="inputField" onChange={handleChange} />
        <input type="password" name="password" placeholder="Password" className="inputField" required onChange={handleChange} />
        <button type="submit" className="submitButton">login</button>
      </form>
      {errorMessage && <InvalidCredentials message={errorMessage} color="red" />}
      <p className="link">
        Create an account? <Link href="/register" className='registerlink'>Register</Link>
      </p>
    </div>
  )
}