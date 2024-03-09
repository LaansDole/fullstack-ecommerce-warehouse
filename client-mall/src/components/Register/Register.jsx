import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { postDataAPI } from "../../api/apiRequest";

import { toast } from "react-hot-toast";

const RegisterComponent = () => {
  const navigate = useNavigate();
  const initialState = {
    username: "",
    password: "",
  };
  const [registerState, setLoginState] = useState(initialState);
  const [confirmPassword, setConfirmPassword] = useState("");

  const rolesMap = ["buyer", "seller"];
  const [roles, setRoles] = useState(rolesMap[0]);
  const [shopName, setShopName] = useState("");
  const { username, password } = registerState;

  const normalCharRegex = /^[A-Za-z0-9._-]*$/;

  const handleChangeInput = e => {
    const { name, value } = e.target;

    // Check if the input matches the allowed characters
    if (!value.match(normalCharRegex)) {
      toast.error("The username must not have strange characters");
      return;
    }

    setLoginState(prevState => ({ ...prevState, [name]: value }));
  };

  const handleRegisterUser = async e => {
    e.preventDefault();

    if (confirmPassword === password) {
      try {
        if (roles === "buyer") {
          const response = await postDataAPI("auth/register", {
            username: username,
            password: password,
            role: roles,
          });
          if (response.status === 200 || response.status === 201) {
            toast.success(`Register Successfully!`);
            // Redirect to the home page
            setTimeout(() => {
              navigate('/');
            }, 1000);
          }
        } else if (roles === "seller") {
          const response = await postDataAPI("auth/register", {
            username: username,
            password: password,
            role: roles,
            shop_name: shopName,
          });
          if (response.status === 200 || response.status === 201) {
            toast.success(`Register Successfully!`);
            // Redirect to the home page
            setTimeout(() => {
              navigate('/');
            }, 1000);
          }
        }
      } catch (error) {
        toast.error(error.response?.data?.error);
      }
    } else {
      toast.error("Your confirmed password does not match");
    }
  };

  return (
    <div className="register_wrapper container">
      <div className="register_container d-flex justify-content-center align-items-center h-100 my-5">
        <div className="register_inner_container d-flex flex-column p-5 text-center">
          <form className="mt-2 mb-5 pb-5" onSubmit={handleRegisterUser}>
            <div className="form_title">
              <h2 className="fw-bold mb-2 text-uppercase">Register</h2>
              <p className="mb-5">Get started with your first account!</p>
            </div>

            <div className="form-floating mb-4">
              <input
                type="text"
                className="form-control"
                id="floatingInput"
                placeholder=""
                name="username"
                value={username}
                onChange={handleChangeInput}
                required
              />
              <label htmlFor="floatingInput">Username</label>
            </div>

            <div className="form-floating mb-4">
              <input
                type="password"
                className="form-control"
                id="floatingPassword"
                placeholder=""
                name="password"
                value={password}
                onChange={handleChangeInput}
                required
              />
              <label htmlFor="floatingPassword">Password</label>
            </div>

            <div className="form-floating mb-4">
              <input
                type="password"
                className="form-control"
                id="floatingPassword"
                placeholder=""
                name="confirmPassword"
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                required
              />
              <label htmlFor="floatingPassword">Confirm Password</label>
            </div>

            <div className="d-flex flex-row mb-4">
              <label htmlFor="radioTitle" className="me-4">
                Roles:{" "}
              </label>
              {rolesMap.map((role, idx) => (
                <div className="form-check me-4" key={idx}>
                  <input
                    className="form-check-input"
                    type="radio"
                    id={`flexRadio${role}`}
                    value={role}
                    onChange={e => setRoles(e.target.value)}
                    checked={roles === role}
                    required
                  />
                  <label
                    className="form-check-label text-capitalize"
                    htmlFor={`flexRadio${role}`}
                  >
                    {role}
                  </label>
                </div>
              ))}
            </div>

            {roles === "seller" && (
              <div className="form-floating mb-4">
                <input
                  type="text"
                  className="form-control"
                  id="floatingInput"
                  placeholder=""
                  name="username"
                  value={shopName}
                  onChange={e => setShopName(e.target.value)}
                  required
                />
                <label htmlFor="floatingInput">Shop Name</label>
              </div>
            )}

            <button
              type="submit"
              className="btn btn-outline-primary w-50 mt-4 px-4 "
            >
              Register
            </button>
          </form>

          <div className="">
            <p className="mb-0">
              Have already an account?
              <Link to="/login" className="fw-bold ms-1">
                Sign In
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RegisterComponent;
