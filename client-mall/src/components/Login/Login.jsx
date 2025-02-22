import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";

import { toast } from "react-hot-toast";
import { postDataAPI } from "../../api/apiRequest";

const LoginComponent = ({ setIsLoggedIn }) => {
  const initialState = {
    username: "",
    password: "",
  };
  const [loginState, setLoginState] = useState(initialState);
  const { username, password } = loginState;

  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || "/";

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

  const handleLoginUser = async e => {
    e.preventDefault();
    try {
      const response = await postDataAPI("auth/login", loginState);
      console.log(response);
      if (response.status === 200 || response.status === 201) {
        localStorage.setItem(
          "userInfo",
          JSON.stringify({
            username: response.data?.username,
            role: response.data?.role,
          }),
        );
        toast.success(
          `Login Successfully! Welcome back ${response.data?.username}`,
        );
        setIsLoggedIn(true);
        navigate(`${from}`, { replace: true });
      }
    } catch (error) {
      toast.error(error.response?.data?.error);
    }
  };

  return (
    <div className="login_wrapper container">
      <div className="login_container d-flex justify-content-center align-items-center h-100">
        <div className="login_inner_container d-flex flex-column p-5 text-center">
          <form className="mt-2 mb-5 pb-5" onSubmit={handleLoginUser}>
            <div className="form_title">
              <h2 className="fw-bold mb-2 text-uppercase">Login</h2>
              <p className="mb-5">Please enter your login and password!</p>
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

            <button
              type="submit"
              className="btn btn-outline-primary w-50 mt-4 px-4 "
            >
              Login
            </button>
          </form>

          <div className="">
            <p className="mb-0">
              Don&apos;t have an account?
              <Link to="/register" className="fw-bold ms-1">
                Sign Up
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

LoginComponent.propTypes = {
  setIsLoggedIn: PropTypes.func.isRequired, // setIsLoggedIn should be a function and is required.
};

export default LoginComponent;
