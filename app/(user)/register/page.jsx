'use client'
import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import '../../styles/register.css'
import Link from "next/link"

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
  const router = useRouter()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleFileChange = (e) => {
    const file = e.target.files[0]
    setFormData({ ...formData, avatar: file })
    
    if (file) {
      const reader = new FileReader()
      reader.onloadend = () => {
        setAvatarPreview(reader.result)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
      
    const formDataToSend = new FormData()
    Object.entries(formData).forEach(([key, value]) => {
      formDataToSend.append(key, value)
    })

    try {
      const response = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        body: formDataToSend,
      })
      
      const result = await response.text()

      if (response.ok) {
        console.log(result)
        router.push("/login")
      } else {
        console.log(result)
        // Handle error - could add state for error message
      }
    } catch (error) {
      console.error('Error:', error)
    }
  }

  return (
    <div className="containereg">
      <h2>Create Your Account</h2>
      <form onSubmit={handleSubmit} className="formreg">
        <div className="inputGroup">
          <div className="inputField">
            <span className="inputLabel">First Name</span>
            <input 
              type="text" 
              name="firstName" 
              placeholder="Enter your first name" 
              required 
              className="inputs" 
              onChange={handleChange} 
              value={formData.firstName}
            />
          </div>
          <div className="inputField">
            <span className="inputLabel">Last Name</span>
            <input 
              type="text" 
              name="lastName" 
              placeholder="Enter your last name" 
              required 
              className="inputs" 
              onChange={handleChange}
              value={formData.lastName}
            />
          </div>
        </div>
        
        <div className="inputField">
          <span className="inputLabel">Email</span>
          <input 
            type="email" 
            name="email" 
            placeholder="Enter your email address" 
            required 
            className="inp" 
            onChange={handleChange}
            value={formData.email}
          />
        </div>
        
        <div className="inputField">
          <span className="inputLabel">Password</span>
          <input 
            type="password" 
            name="password" 
            placeholder="Create a password" 
            required 
            className="inp" 
            onChange={handleChange}
            value={formData.password}
          />
        </div>
        
        <div className="inputGroup">
          <div className="inputField">
            <span className="inputLabel">Date of Birth</span>
            <input 
              type="date" 
              name="dateOfBirth" 
              required 
              className="inputs" 
              onChange={handleChange}
              value={formData.dateOfBirth}
            />
          </div>
          <div className="inputField">
            <span className="inputLabel">Nickname (Optional)</span>
            <input 
              type="text" 
              name="nickname" 
              placeholder="Choose a nickname" 
              className="inputs" 
              onChange={handleChange}
              value={formData.nickname}
            />
          </div>
        </div>
        
        <div className="inputsavatar">
          <div className="avatarPreview">
            {avatarPreview ? (
              <img src={avatarPreview} alt="Avatar preview" />
            ) : (
              <span>No image</span>
            )}
          </div>
          <div className="avatarUpload">
            <input 
              type="file" 
              name="avatar" 
              accept="image/*" 
              onChange={handleFileChange} 
              id="avatarInput"
            />
            <label htmlFor="avatarInput" className="avatarUploadBtn">
              Choose Profile Picture
            </label>
          </div>
        </div>
        
        <div className="inputField">
          <span className="inputLabel">About Me (Optional)</span>
          <textarea 
            name="aboutMe" 
            placeholder="Tell us a bit about yourself" 
            className="inp" 
            onChange={handleChange}
            value={formData.aboutMe}
            rows="3"
          ></textarea>
        </div>
        
        <button type="submit" className="submitButtonreg">Create Account</button>
        
        <div className="linkreg">
          Already have an account? <Link href="/login" className="loginregisterlink">Sign In</Link>
        </div>
      </form>
    </div>
  )
}