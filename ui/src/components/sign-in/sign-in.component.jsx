import React, {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';

import Input from '../input/input.component';
import Button from '../button/button.component';
import ErrorMessage from '../messages/error/error-message.component';
import signInMutation from '../../graphql/mutations/sign-in-mutation';

import './sign-in.styles.scss';

const SignIn = ({signInCallback}) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState(null);
  const [signIn, {loading}] = useMutation(signInMutation, {
    onCompleted: (data) => {
      localStorage.setItem('authenticationToken', data.signIn.token);
      localStorage.setItem('userId', data.signIn.user.id);
      localStorage.setItem('username', data.signIn.user.username);
      signInCallback(localStorage.getItem('authenticationToken'));
    },
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
    }
  });

  const onSubmit = (event) => {
    event.preventDefault();
    signIn({variables: {email: email, password: password}});
  }

  let submitButton;
  loading ? submitButton = <Button disabled={true}> Sign in </Button> :
    submitButton = <Button type='submit'> Sign in </Button>;

  let errorMessage;
  errors ? errorMessage = errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>;
  }) : errorMessage = '';

  return (
    <div className='sign-in'>
      {errorMessage}
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
