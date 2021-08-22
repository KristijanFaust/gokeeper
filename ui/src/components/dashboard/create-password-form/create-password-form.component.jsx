import {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';
import {useHistory} from 'react-router-dom';

import Input from '../../input/input.component';
import Button from '../../button/button.component';
import ErrorMessage from '../../messages/error/error-message.component';
import createPasswordMutation from '../../../graphql/mutations/create-password-mutation';

import './create-password-form.component.scss'

const CreatePasswordForm = ({createPasswordsCallback, authenticationExpiredCallback}) => {
  let history = useHistory();
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState(null);
  const [createPassword, {loading}] = useMutation(createPasswordMutation, {
    onCompleted: (data) => {
      setName('');
      setPassword('');
      createPasswordsCallback(data.createPassword);
    },
    onError: (response) => {
      setErrors(response.graphQLErrors?.map(error => error.message));
      if (!response.graphQLErrors?.length) {
        localStorage.clear();
        authenticationExpiredCallback('');
        history.push('/sign-in', {authenticationExpired: true});
      }
    }
  });

  const onSubmit = (event) => {
    event.preventDefault();
    createPassword({variables: {userId: localStorage.getItem('userId'), name: name, password: password}});
  }

  const submitButton = loading ? <Button disabled={true}>Save password</Button>
    : <Button type='submit'>Save password</Button>;

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  return (
    <div>
      {errorMessage}
      <form className='create-password-form' onSubmit={onSubmit}>
        <Input
          name='name' type='text' label='name' required
          onChange={event => setName(event.target.value)} value={name}
        />
        <Input
          name='password' type='text' label='password' required
          onChange={event => setPassword(event.target.value)} value={password}
        />
        {submitButton}
      </form>
    </div>
  );
};

export default CreatePasswordForm;