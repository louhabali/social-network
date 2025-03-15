'use client'
import { useState, useRef, useEffect } from "react"
import { useRouter } from "next/navigation"
import '../../styles/register.css'
import Link from "next/link"
<<<<<<< HEAD

=======
import InvalidCredentials from '@/components/invalidCredentials'
>>>>>>> alouhab
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
<<<<<<< HEAD
  
  const [avatarPreview, setAvatarPreview] = useState(null)
=======
  const [avatarPreview, setAvatarPreview] = useState(null)
  const [errorMessage, setErrorMessage] = useState("");
>>>>>>> alouhab
  const router = useRouter()

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleFileChange = (e) => {
<<<<<<< HEAD
    const file = e.target.files[0]
    setFormData({ ...formData, avatar: file })
    
    if (file) {
=======
    setFormData({ ...formData, avatar: e.target.files[0] })
    if ( e.target.files[0]) {
>>>>>>> alouhab
      const reader = new FileReader()
      reader.onloadend = () => {
        setAvatarPreview(reader.result)
      }
<<<<<<< HEAD
      reader.readAsDataURL(file)
=======
      reader.readAsDataURL( e.target.files[0])
>>>>>>> alouhab
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
<<<<<<< HEAD
      
=======
  
    if (formData.email.trim().length==0 || formData.password.trim().length==0 || formData.firstName.trim().length==0 || formData.lastName.trim().length==0 || formData.dateOfBirth.trim().length==0) {
      setErrorMessage("Please fill all the required fields")
      return
    }

>>>>>>> alouhab
    const formDataToSend = new FormData()
    Object.entries(formData).forEach(([key, value]) => {
      formDataToSend.append(key, value)
    })

    try {
      const response = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        body: formDataToSend,
      })
<<<<<<< HEAD
      
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
=======

      if (!response.ok) {
        const result = await response.text()
        setErrorMessage(result)
      } else {
        router.push("/login")
      }
    } catch (error) {
      console.error('Error', error)
>>>>>>> alouhab
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
<<<<<<< HEAD
=======
        {errorMessage && <InvalidCredentials message={errorMessage} color="red" />}
>>>>>>> alouhab
        
        <div className="linkreg">
          Already have an account? <Link href="/login" className="loginregisterlink">Sign In</Link>
        </div>
      </form>
    </div>
  )
}