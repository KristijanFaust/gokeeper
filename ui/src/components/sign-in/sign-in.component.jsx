import {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';
import { useLocation } from 'react-router-dom';
import { useEffect } from 'react';

import Input from '../input/input.component';
import Button from '../button/button.component';
import ErrorMessage from '../messages/error/error-message.component';
import NotificationMessage from '../messages/notification/notification-message.component';
import signInMutation from '../../graphql/mutations/sign-in-mutation';

import './sign-in.styles.scss';

const SignIn = ({signInCallback}) => {
  const location = useLocation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState(null);
  const [notifications, setNotifications] = useState(null);
  const [signIn, {loading}] = useMutation(signInMutation, {
    onCompleted: (data) => {
      localStorage.setItem('authenticationToken', data.signIn.token);
      localStorage.setItem('userId', data.signIn.user.id);
      localStorage.setItem('username', data.signIn.user.username);
      signInCallback(localStorage.getItem('authenticationToken'));
    },
    onError: (response) => {
      setErrors(response.graphQLErrors?.map(error => error.message));
    }
  });

  useEffect(() => {
    if (location.state?.email) {
      setEmail(location.state.email);
      setNotifications(['You are signed up!']);
    }
  }, [location]);

  const onSubmit = (event) => {
    event.preventDefault();
    signIn({variables: {email: email, password: password}});
  }

  const submitButton = loading ? <Button disabled={true}> Sign in </Button> : <Button type='submit'> Sign in </Button>;

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  const notificationMessage = notifications ? notifications.map((notification, index) => {
    return <NotificationMessage key={index}>{notification}</NotificationMessage>
  }) : null;

  return (
    <div className='sign-in'>
      {errorMessage}
      {notificationMessage}
      <form onSubmit={onSubmit}>
        <Input
          name='email' type='email' label='email' required
          onChange={event => setEmail(event.target.value)} value={email}
        />
        <Input
          name='password' type='password' label='password' minLength={8} required
          onChange={event => setPassword(event.target.value)} value={password}
        />
        {submitButton}
      </form>
    </div>
  );
};

export default SignIn;
