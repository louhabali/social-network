'use client';
import styles from './login.module.css'
import { useRouter } from "next/navigation"
import { useState } from "react"
import Link from "next/link"
import InvalidCredentials from '@/components/invalidCredentials'

export default function LoginPage() {
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
    if (formData.email.trim().length == 0 || formData.password.trim().length == 0) {
      setErrorMessage("Please fill all the required fields")
      return
    }

    if (!formData.email.match(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/)) {
      setErrorMessage("Please enter a valid email");
      return
    } 
    if (formData.password.length < 8) {
      setErrorMessage("Password must be at least 8 characters long");
      return
    }

    if (!formData.password.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/)) {
      setErrorMessage("Password must contain at least one uppercase letter, one lowercase letter, and one number");
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
        router.push("/");
      } else {
        setErrorMessage(result);
      }
    } catch (error) {
      console.error('Error', error)
    }
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.container}>
        <h2>Welcome Back</h2>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.inputField}>
            <span className={styles.inputLabel}>Email</span>
            <input
              type="email"
              name="email"
              placeholder="Enter your email"
              required
              className={styles.input}
              onChange={handleChange}
              value={formData.email}
            />
          </div>

          <div className={styles.inputField}>
            <span className={styles.inputLabel}>Password</span>
            <input
              type="password"
              name="password"
              placeholder="Enter your password"
              className={styles.input}
              required
              onChange={handleChange}
              value={formData.password}
            />
          </div>

          <button type="submit" className={styles.submitButton}>Sign In</button>
          {errorMessage && <InvalidCredentials message={errorMessage} color="red" />}
        </form>

        <p className={styles.link}>
          Don't have an account? <Link href="/register" className={styles.registerLink}>Create Account</Link>
        </p>
      </div>
    </div>
  )
}