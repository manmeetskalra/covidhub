import React, { useState, createContext } from "react";

export const AuthContext = createContext();

export default ({ children }) => {
  const [email, setEmail] = useState(null);
  const [password, setPassword] = useState(null);
  const [phoneNumber, setPhoneNumber] = useState(null);

  return (
    <div>
      <AuthContext.Provider
        value={{
          email,
          setEmail,
          password,
          setPassword,
          phoneNumber,
          setPhoneNumber
        }}
      >
        {children}
      </AuthContext.Provider>
    </div>
  );
};
