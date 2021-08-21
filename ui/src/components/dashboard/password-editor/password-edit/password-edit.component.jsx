import {faEye, faSave, faTrash} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';

import FormInput from '../../../input/input.component';
import Button from '../../../button/button.component';
import updatePasswordMutation from '../../../../graphql/mutations/update-password-mutation';

import './password-edit.component.scss';
import ErrorMessage from '../../../messages/error/error-message.component';

const PasswordEdit = ({passwordEntry}) => {
  const [name, setName] = useState(passwordEntry.name)
  const [password, setPassword] = useState(passwordEntry.password)
  const [showPassword, setShowPassword] = useState(false)
  const [errors, setErrors] = useState(null);
  const [updatePassword, {loading: updateLoading}] = useMutation(updatePasswordMutation, {
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
    }
  });

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const updatePasswordHandler = () => {
    updatePassword({variables: {passwordId: passwordEntry.id, name: name, password: password}});
  };

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  return (
    <div className='password-edit-container'>
      {errorMessage}
      <div className='password-edit'>
        <FormInput type='text' onChange={event => setName(event.target.value)} value={name} />
        <FormInput type={showPassword ? 'text' : 'password'} onChange={event => setPassword(event.target.value)} value={password} />
        <Button onClick={() => togglePasswordVisibility()}>
          <FontAwesomeIcon icon={faEye} />
        </Button>
        <Button onClick={() => updatePasswordHandler()} disabled={updateLoading ?? false}>
          <FontAwesomeIcon icon={faSave} />
        </Button>
        <Button>
          <FontAwesomeIcon icon={faTrash} />
        </Button>
      </div>
    </div>
  );
};

export default PasswordEdit;