'use client'
import { useState } from "react"
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

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleFileChange = (e) => {
    setFormData({ ...formData, avatar: e.target.files[0] })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    const response = await fetch("http://localhost:8080/api/register", {
      method: "POST",
      body: JSON.stringify(formData),
    })

    const result = await response.text()
    console.log(result)
  }

  return (
    <div className="containereg">
      <form onSubmit={handleSubmit} className="formreg">
        <div className="inputsd">
        <input type="email" name="email" placeholder="Email" required className="inputs" onChange={handleChange} />
        <input type="password" name="password" placeholder="Password" required className="inputs" onChange={handleChange} />
        </div>
        <div className="inputsd">
        <input type="text" name="firstName" placeholder="First Name" required className="inputs" onChange={handleChange} />
        <input type="text" name="lastName" placeholder="Last Name" required className="inputs" onChange={handleChange} />
        </div>
        <div className="inputsd">
        <input type="date" name="dateOfBirth" required className="inputs" onChange={handleChange} /> 
        <input type="text" name="nickname" placeholder="Nickname (Optional)"className="inputs" onChange={handleChange} />
        </div>
        <div className="inputsavatar">
        <input type="file" name="avatar" accept="image/*" className="inp" onChange={handleFileChange} />
        <textarea name="aboutMe" placeholder="About Me (Optional)"className="inp" onChange={handleChange}></textarea>
        </div>
        <div className="inputsbtn">
        <button type="submit" className="submitButtonreg">Register</button>
        <p className="linkreg">
        already have an acount? <Link href="/login" className='loginregisterlink'>login</Link>
        </p>
        </div>
      </form>
      
      
    </div>
  )
}
