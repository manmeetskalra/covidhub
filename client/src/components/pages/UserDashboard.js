import React, { useState, useRef, useEffect, useContext } from "react";
import Button from "@material-ui/core/Button";
import CssBaseline from "@material-ui/core/CssBaseline";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import Container from "@material-ui/core/Container";
import {
  ValidatorForm,
  TextValidator,
  SelectValidator
} from "react-material-ui-form-validator";
import { useHistory } from "react-router-dom";
//import ApiService from './../../services/ApiService';
import Snackbar from "@material-ui/core/Snackbar";
import MuiAlert from "@material-ui/lab/Alert";
import ApiService from "./../../services/ApiService";
import { AuthContext } from "./../../context/ContextApi";
import InputLabel from "@material-ui/core/InputLabel";
import MenuItem from "@material-ui/core/MenuItem";
import FormHelperText from "@material-ui/core/FormHelperText";
import FormControl from "@material-ui/core/FormControl";
import ListItemText from "@material-ui/core/ListItemText";
import Select from "@material-ui/core/Select";
import Checkbox from "@material-ui/core/Checkbox";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import AppBar from "@material-ui/core/AppBar";
import TextField from "@material-ui/core/TextField";
import Link from "@material-ui/core/Link";

const countryOptions = [
  "India",
  "Sri Lanka",
  "China",
  "Singapore",
  "UK",
  "Canada",
  "Bangladesh",
  "Pakistan",
  "America"
];

const informationType = ["Confirmed", "Recovered", "Deaths"];

const frequencyOptions = ["Once", "Twice", "Thrice"];

function Alert(props) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

const useStyles = makeStyles(theme => ({
  paper: {
    marginTop: "30%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    borderColor: theme.palette.primary.dark,
    border: "2px solid",
    padding: "60px"
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main
  },
  form: {
    width: "100%",
    marginTop: theme.spacing(1)
  },
  submit: {
    margin: theme.spacing(3, 0, 2)
  }
}));

export default function UserDashboard() {
  var { email, phoneNumber } = useContext(AuthContext);
  /*
  Here email and password are held in a global state. Use these variables to authenticate by using the
  /authenticate endpoint on the server. If the user is not authenticated, don't load any componenet
  */

  const classes = useStyles();
  const formRef = useRef();
  const history = useHistory();
  const [open, setOpen] = useState(false);
  const [displayEmail, setDisplayEmail] = useState(true);
  const [country, setCountry] = useState([]);
  const [infoType, setInfoType] = useState("");
  const [frequency, setFrequency] = useState("");
  const [user, setUser] = useState({
    servicType: "",
    country: "",
    informationType: "",
    frequency: ""
  });

  useState(() => {
    console.log(email, phoneNumber);
    if (email == null) email = localStorage.getItem("email");
    if (phoneNumber == null) phoneNumber = localStorage.getItem("phoneNumber");
  }, []);

  const handleSubmit = e => {
    e.preventDefault();
    let freq = 86400;
    freq /= frequency == "Once" ? 1 : frequency == "Twice" ? 2 : 3;
    let countriesMentioned = country.map(c => c.toLowerCase());
    let info = {
      countries: countriesMentioned.join("|"),
      frequency: freq
    };
    if (displayEmail) {
      info["email"] = email;
      info["type"] = infoType.toLowerCase();
      ApiService.subscribeEmail(info).then(data => {
        console.log(data);
      });
    } else {
      info["phoneNumber"] = phoneNumber;
      ApiService.subscribeText(info);
    }
  };

  const handleSubscribe = e => {
    e.preventDefault();
    let freq = 86400;
    freq /= frequency == "Once" ? 1 : frequency == "Twice" ? 2 : 3;
    let countriesMentioned = country.map(c => c.toLowerCase());
    let info = {
      countries: countriesMentioned.join("|"),
      frequency: freq
    };

    if (displayEmail) {
      info["email"] = email;
      info["type"] = infoType.toLowerCase();
      ApiService.subscribeNowEmail(info).then(data => {
        console.log(data);
      });
    } else {
      info["phoneNumber"] = phoneNumber;
      ApiService.subscribeNowText(info).then(data => {
        console.log(data);
      });
    }
  };

  const handleUnsubscribe = e => {
    e.preventDefault();
    if (displayEmail) {
      ApiService.unsubscribeEmail({
        email
      });
    } else {
      ApiService.unsubscribeText({
        phoneNumber
      });
    }
  };

  const handleChange = e => {
    setUser({ ...user, [e.target.name]: e.target.value });
  };

  const setNotification = () => {
    setOpen(true);
  };

  const handleClose = (event, reason) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
  };

  const handleTabChange = (event, newValue) => {
    // newValue is 1 for text and 0 for email
    if (newValue == 0 && !displayEmail) {
      setDisplayEmail(true);
    } else if (newValue == 1 && displayEmail) {
      setDisplayEmail(false);
    }
  };

  const handleCountryChange = event => {
    setCountry(event.target.value);
  };

  const handleInfoChange = event => {
    setInfoType(event.target.value);
  };

  const handleFrequencyChange = event => {
    console.log(event.target.value);
    setFrequency(event.target.value);
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <div className={classes.paper}>
        <Typography component="h1" variant="h5">
          Service Information
        </Typography>

        <AppBar position="static">
          <Tabs onChange={handleTabChange} aria-label="simple tabs example">
            <Tab label="E-mail" />
            <Tab label="Text" />
          </Tabs>
        </AppBar>

        <TextField
          id="service-type-email"
          defaultValue={email}
          disabled
          margin="normal"
          fullWidth
          label="Email"
          variant="outlined"
          style={{
            marginTop: 25
          }}
        />
        <TextField
          id="service-type-phoneNumber"
          defaultValue={phoneNumber}
          disabled
          margin="normal"
          fullWidth
          label="Phone Number"
          variant="outlined"
        />

        <ValidatorForm
          className={classes.form}
          ref={formRef}
          onSubmit={handleSubmit}
          onError={errors => console.log(errors)}
        >
          <InputLabel
            id="country-input"
            style={{
              marginTop: 25
            }}
          >
            Country(s)
          </InputLabel>
          <Select
            labelId="multiple-countires"
            id="countries-selection"
            multiple
            value={country}
            onChange={handleCountryChange}
            renderValue={selected => selected.join(", ")}
          >
            {countryOptions.map(name => (
              <MenuItem key={name} value={name}>
                <Checkbox checked={country.indexOf(name) > -1} />
                <ListItemText primary={name} />
              </MenuItem>
            ))}
          </Select>
          {displayEmail && (
            <SelectValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              value={infoType}
              name="informationType"
              label="Information Type"
              onChange={handleInfoChange}
              validators={["required"]}
              errorMessages={["this field is required"]}
            >
              {informationType.map(type => (
                <MenuItem key={type} value={type}>
                  {type}
                </MenuItem>
              ))}
            </SelectValidator>
          )}
          <SelectValidator
            variant="outlined"
            margin="normal"
            required
            fullWidth
            value={frequency}
            name="frequency"
            label="Frequency"
            onChange={handleFrequencyChange}
            validators={["required"]}
            errorMessages={["this field is required"]}
          >
            {frequencyOptions.map(freq => (
              <MenuItem key={freq} value={freq}>
                {freq}
              </MenuItem>
            ))}
          </SelectValidator>

          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
            onClick={handleSubmit}
          >
            Subscribe
          </Button>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
            onClick={handleSubscribe}
          >
            Get instant notification
          </Button>
          <Typography variant="caption" align="center">
            If you wish to unsubscribe, click{" "}
            <Link href="#" onClick={handleUnsubscribe}>
              here
            </Link>
          </Typography>
        </ValidatorForm>
      </div>
    </Container>
  );
}
