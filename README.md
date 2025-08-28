# API Planning

## Project Goal
Create an API for a school management system that administrative staff can use to manage students, teachers, and other staff members.

## Key Requirements
- Addition of student/teacher/staff/exec entry  
- Modification of student/teacher/staff/exec entry  
- Delete student/teacher/staff/exec entry  
- Get list of all students/teachers/staff/execs  
- Authentication: Login, Logout  
- Bulk Modifications: students/teachers/staff/execs  

### Class Management
- Total count of a class with class teacher  
- List of all students in a class with class teacher  

### Security and Rate Limiting
- Rate limit the application  
- Password reset mechanisms (forgot password, update password)  
- Deactivate user  

---

## Executives

- `GET /execs`: Get list of executives  
- `POST /execs`: Add a new executive  
- `PATCH /execs`: Modify multiple executives  
- `GET /execs/{id}`: Get a specific executive  
- `PATCH /execs/{id}`: Modify a specific executive  
- `DELETE /execs/{id}`: Delete a specific executive  
- `POST /execs/login`: Login  
- `POST /execs/logout`: Logout  
- `POST /execs/forgotpassword`: Forgot password  
- `POST /execs/resetpassword/reset/{resetcode}`: Reset password  

---

## Students

- `GET /students`: Get list of students  
- `POST /students`: Add a new student  
- `PATCH /students`: Modify multiple students  
- `DELETE /students`: Delete multiple students  
- `GET /students/{id}`: Get a specific student  
- `PUT /students/{id}`: Update a specific student  
- `PATCH /students/{id}`: Modify a specific student  
- `DELETE /students/{id}`: Delete a specific student  

---

## Teachers

- `GET /teachers`: Get list of teachers  
- `POST /teachers`: Add a new teacher  
- `PATCH /teachers`: Modify multiple teachers  
- `DELETE /teachers`: Delete multiple teachers  
- `GET /teachers/{id}`: Get a specific teacher  
- `PUT /teachers/{id}`: Update a specific teacher  
- `PATCH /teachers/{id}`: Modify a specific teacher  
- `DELETE /teachers/{id}`: Delete a specific teacher  
- `GET /teachers/{id}/students`: Get students of a specific teacher  
- `GET /teachers/{id}/studentcount`: Get student count for a specific teacher  


## Best Practices
- Modularity  
- Documentation  
- Error Handling  
- Security  
- Testing  

## Common Pitfalls
- Overcomplicating the API  
- Ignoring Security  
- Poor Documentation  
- Inadequate Testing  
