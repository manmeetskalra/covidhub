import React, {useState, useRef, useEffect} from 'react';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import { ValidatorForm, TextValidator} from 'react-material-ui-form-validator';
import { useHistory } from "react-router-dom";
import ApiService from './../../services/ApiService';
import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert from '@material-ui/lab/Alert';

function Alert(props) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: '30%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    borderColor: theme.palette.primary.dark,
    border: '2px solid',
    padding: '60px'
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: '100%',
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
     
}));

export default function Register(props, context) {

  const classes = useStyles();
  const formRef = useRef();
  const history = useHistory();
  const [open, setOpen] = useState(false);
  const [user, setUser] = useState({email: "", firstname: "", lastname: "", phonenumber: "+", password: ""});


  const handleSubmit = e => {
    e.preventDefault();
    
    ApiService.register(user).then(data => { 
      const { response } = data;
      console.log(response);

      if(response.msgError){
        setNotification();
      }else {
        history.push('/');
      }
    })
  }

  const handleChange = e => {
    setUser({...user,[e.target.name] : e.target.value});
  }

  const setNotification = () => {
    setOpen(true);
  };

  const handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    setOpen(false);
  };


  useEffect(() => {
 
    ValidatorForm.addValidationRule('minLength', value => {
      if (value.length < 3) {
          return false;
      }
      return true;
    });

    ValidatorForm.addValidationRule('maxLength', value => {
      if (value.length > 15) {
          return false;
      }
      return true;
    });

    ValidatorForm.addValidationRule('validNumber', value => {
      //NOTE: Twilio recommends E.164 formatting: +14155552671
      if (value.length === 12 && value[0] === '+' && !isNaN(value.substring(1,value.length))) {
          return true;
      }
      return false;
    });

    ValidatorForm.addValidationRule('validPassword', value => {
      console.log(value);
      if (/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/.test(value)) {
          return true;
      }
      return false;
    });
  });

  return (
    <Container component="main" maxWidth="xs" >
        <CssBaseline />
        <div className={classes.paper}>
          <Snackbar open={open} autoHideDuration={6000} onClose={handleClose}>
            <Alert onClose={handleClose} severity="error">
              Email already exists
            </Alert>
          </Snackbar>
          <Typography component="h1" variant="h5">
            Register
          </Typography>
          <ValidatorForm
            className={classes.form}
            ref={formRef}
            onSubmit={handleSubmit}
            onError={errors => console.log(errors)}
          >
            <TextValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="firstname"
              label="First Name"
              autoComplete="First Name"
              onChange={handleChange}
              value={user.firstname}
              validators={['required', 'minLength', 'maxLength']}
              errorMessages={['this field is required', 'must be greater than 3 characters', 'must be less than 16 characters']}
            />
            <TextValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="lastname"
              label="Last Name"
              autoComplete="Last Name"
              onChange={handleChange}
              value={user.lastname}
              validators={['required', 'minLength', 'maxLength']}
              errorMessages={['this field is required', 'must be greater than 3 characters', 'must be less than 16 characters']}
            />
            <TextValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="phonenumber"
              label="Phone Number"
              autoComplete="Phone Number"
              id='phonenumberForm'
              onChange={handleChange}
              value={user.phonenumber}
              validators={['required', 'validNumber']}
              errorMessages={['this field is required', 'must be a + with an 11 digit number']}
            />
            <TextValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              label="Email"
              onChange={handleChange}
              name="email"
              value={user.email}
              validators={['required', 'isEmail']}
              errorMessages={['this field is required', 'email is not valid']}
            />
            <TextValidator
              variant="outlined"
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"   
              onChange={handleChange}
              value={user.password}
              validators={['required', 'validPassword']}
              errorMessages={['this field is required', 'minimum eight characters, at least one letter, one number and one special character']}
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              className={classes.submit}
            >
              Register
            </Button>
          </ValidatorForm>
      </div>
    </Container>
  );
}