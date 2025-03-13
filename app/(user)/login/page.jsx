'use client';
import '../../styles/login.css'
import { useRouter } from "next/navigation"
import { useState } from "react"
import Link from "next/link"

export default function LoginPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })

  const router = useRouter()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    const response = await fetch("http://localhost:8080/api/login", {
      method: "POST",
      body: JSON.stringify(formData),
      credentials: "include",
    })

    const result = await response.text()

    if (response.ok) {
      console.log(result);
      router.push("/")
    } else {
      console.log(result);
    }
  }

  return (
    <div className="container">
      <h2>Welcome Back</h2>
      <form onSubmit={handleSubmit} className="form">
        <div className="inputField">
          <span className="inputLabel">Email</span>
          <input
            type="email"
            name="email"
            placeholder="Enter your email"
            required
            className="input"
            onChange={handleChange}
            value={formData.email}
          />
        </div>

        <div className="inputField">
          <span className="inputLabel">Password</span>
          <input
            type="password"
            name="password"
            placeholder="Enter your password"
            className="input"
            required
            onChange={handleChange}
            value={formData.password}
          />
        </div>

        <div className="forgotPassword">
          <Link href="/forgot-password" className="forgotPasswordLink">
            Forgot password?
          </Link>
        </div>

        <button type="submit" className="submitButton">Sign In</button>

        <div className="formDivider">
          <span>or</span>
        </div>

        <div className="socialLogin">
          <button type="button" className="socialButton">G</button>
          <button type="button" className="socialButton">f</button>
          <button type="button" className="socialButton">in</button>
        </div>
      </form>

      <p className="link">
        Don't have an account? <Link href="/register" className="registerLink">Create Account</Link>
      </p>
    </div>
  )
}