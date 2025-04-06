'use client'
import { useState } from "react"
import { useRouter } from "next/navigation"
import styles from "./register.module.css"
import Link from "next/link"
import InvalidCredentials from '@/components/invalidCredentials'

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    dateOfBirth: "",
    avatar: null,
    nickname: "",
    aboutMe: ""
  })
  const [avatarPreview, setAvatarPreview] = useState(null)
  const [errorMessage, setErrorMessage] = useState("");
  const router = useRouter()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 1024 * 1024) {
        setErrorMessage("Profile image must be less than 2MB");
        return;
      }
      const allowedTypes = ['image/jpeg', 'image/png'];
      if (!allowedTypes.includes(file.type)) {
        setErrorMessage("Only JPG or PNG images are allowed");
        return;
      }
      setAvatarPreview(URL.createObjectURL(file));
      setFormData({ ...formData, avatar: e.target.files[0] })
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (formData.email.trim().length == 0 || formData.password.trim().length == 0 || formData.firstName.trim().length == 0 || formData.lastName.trim().length == 0 || formData.dateOfBirth.trim().length == 0) {
      setErrorMessage("Please fill all the required fields")
      return
    }

    if (!formData.email.match(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/)) {
      setErrorMessage("Please enter a valid email" );
      return
    }
    if (!formData.password.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/)) {
      setErrorMessage("Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number");
      return
    }
    if (!formData.firstName.match(/^[a-zA-Z]{1,10}$/) || !formData.lastName.match(/^[a-zA-Z]{1,10}$/)) {
      setErrorMessage("First and last name must contain only letters and be less than 10 characters each");
      return;
    }
    if (formData.nickname && !formData.nickname.match(/^[a-zA-Z0-9]{1,10}$/)) {
      setErrorMessage("Nickname must contain only letters and numbers and be less than 10 characters long");
      return
    }
    if (formData.aboutMe && formData.aboutMe.length > 100) {
      setErrorMessage("About me must be less than 100 characters long");
      return
    }
    
    if (formData.dateOfBirth) {
      const dateOfBirth = new Date(formData.dateOfBirth) 
      const today = new Date()
      let age = today.getFullYear() - dateOfBirth.getFullYear()
      if (today.getMonth() - dateOfBirth.getMonth() < 0 || (dateOfBirth.getMonth() === today.getMonth() && today.getDate() < dateOfBirth.getDate())) {
        age--
      }
      if (age < 5) {
        setErrorMessage("You must be at least 5 years old to register" );
        return
      }
    }

    const formDataToSend = new FormData()
    Object.entries(formData).forEach(([key, value]) => {
      formDataToSend.append(key, value)
    })

    try {
      const response = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        body: formDataToSend,
      })

      if (!response.ok) {
        const result = await response.text()
        setErrorMessage(result)
      } else {
        router.push("/login")
      }
    } catch (error) {
      console.error('Error', error)
    }
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.containereg}>
        <h2>Create Your Account</h2>
        <form onSubmit={handleSubmit} className={styles.formreg}>
          <div className={styles.inputGroup}>
            <div className={styles.inputField}>
              <span className={styles.inputLabel}>First Name</span>
              <input type="text" name="firstName" placeholder="Enter your first name" required className={styles.inputs} onChange={handleChange} value={formData.firstName} />
            </div>
            <div className={styles.inputField}>
              <span className={styles.inputLabel}>Last Name</span>
              <input type="text" name="lastName" placeholder="Enter your last name" required className={styles.inputs} onChange={handleChange} value={formData.lastName} />
            </div>
          </div>

          <div className={styles.inputField}>
            <span className={styles.inputLabel}>Email</span>
            <input type="email" name="email" placeholder="Enter your email address" required className={styles.inp} onChange={handleChange} value={formData.email} />
          </div>

          <div className={styles.inputField}>
            <span className={styles.inputLabel}>Password</span>
            <input type="password" name="password" placeholder="Create a password" required className={styles.inp} onChange={handleChange} value={formData.password} />
          </div>

          <div className={styles.inputGroup}>
            <div className={styles.inputField}>
              <span className={styles.inputLabel}>Date of Birth</span>
              <input type="date" name="dateOfBirth" required className={styles.inputs} onChange={handleChange} value={formData.dateOfBirth} max={new Date(new Date().setFullYear(new Date().getFullYear() - 5)).toISOString().split("T")[0]}/>
            </div>
            <div className={styles.inputField}>
              <span className={styles.inputLabel}>Nickname (Optional)</span>
              <input type="text" name="nickname" placeholder="Choose a nickname" className={styles.inputs} onChange={handleChange} value={formData.nickname} />
            </div>
          </div>

          <div className={styles.inputsavatar}>
            <div className={styles.avatarPreview}>
              {avatarPreview ? <img src={avatarPreview} alt="Avatar preview" /> : <span>No image</span>}
            </div>
            <div className={styles.avatarUpload}>
              <input type="file" name="avatar" accept="image/*" onChange={handleFileChange} id="avatarInput" style={{ display: "none" }} />
              <label htmlFor="avatarInput" className={styles.avatarUploadBtn}>
                {avatarPreview ? "Change Image" : "Choose Profile Picture"}
              </label>
              {avatarPreview && (
                <button type="button" className={styles.avatarUploadBtn} style={{ backgroundColor: "red" }} onClick={() => { setAvatarPreview(null); setFormData({ ...formData, avatar: null });document.getElementById("avatarInput").value = "";}}>
                  Cancel
                </button>
              )}
            </div>
          </div>

          <div className={styles.inputField}>
            <span className={styles.inputLabel}>About Me (Optional)</span>
            <textarea name="aboutMe" placeholder="Tell us a bit about yourself" className={styles.inp} onChange={handleChange} value={formData.aboutMe} rows="3"></textarea>
          </div>

          <button type="submit" className={styles.submitButtonreg}>Create Account</button>
          {errorMessage && <InvalidCredentials message={errorMessage} color="red" />}

          <div className={styles.linkreg}>
            Already have an account? <Link href="/login" className={styles.loginregisterlink}>Sign In</Link>
          </div>
        </form>
      </div>
    </div>
  )
}
