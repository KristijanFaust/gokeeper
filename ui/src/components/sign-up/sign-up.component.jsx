import {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';
import { useHistory } from 'react-router-dom';

import Input from '../input/input.component';
import Button from '../button/button.component';
import ErrorMessage from '../messages/error/error-message.component';
import signUpMutation from '../../graphql/mutations/sign-up-mutation';

import './sign-up.styles.scss';

const SignUp = () => {
  let history = useHistory();
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState(null);
  const [signUp, {loading}] = useMutation(signUpMutation, {
    onCompleted: (data) => {
      history.push('/sign-in', {email: data.signUp.email});
    },
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
    }
  });

  const onSubmit = (event) => {
    event.preventDefault();
    signUp({variables: {email: email, username: username, password: password}});
  }

  const submitButton = loading ? <Button disabled={true}> Sign up </Button> : <Button type='submit'> Sign up </Button>;

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  return (
    <div className='sign-up'>
      {errorMessage}
      <form onSubmit={onSubmit}>
        <Input
          name='email' type='email' label='email' required
          onChange={event => setEmail(event.target.value)} value={email}
        />
        <Input
          name='username' type='text' label='username' required
          onChange={event => setUsername(event.target.value)} value={username}
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

export default SignUp;
