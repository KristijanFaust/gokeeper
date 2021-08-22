import {faEye, faSave, faTrash} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {useState} from 'react';
import {useMutation} from '@apollo/react-hooks';
import {useHistory} from 'react-router-dom';

import FormInput from '../../../input/input.component';
import Button from '../../../button/button.component';
import ErrorMessage from '../../../messages/error/error-message.component';
import updatePasswordMutation from '../../../../graphql/mutations/update-password-mutation';
import deletePasswordMutation from '../../../../graphql/mutations/delete-password-mutation';

import './password-edit.component.scss';

const PasswordEdit = ({passwordEntry, authenticationExpiredCallback}) => {
  let history = useHistory();
  const [name, setName] = useState(passwordEntry.name)
  const [password, setPassword] = useState(passwordEntry.password)
  const [showPassword, setShowPassword] = useState(false)
  const [errors, setErrors] = useState(null);
  const [deleted, setDeleted] = useState(false)
  const [updatePassword, {loading: updateLoading}] = useMutation(updatePasswordMutation, {
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
      if (!response.graphQLErrors.length) {
        localStorage.clear();
        authenticationExpiredCallback('')
        history.push('/sign-in', {authenticationExpired: true});
      }
    }
  });
  const [deletePassword, {loading: deleteLoading}] = useMutation(deletePasswordMutation, {
    onCompleted: () => {
      setDeleted(true)
    },
    onError: (response) => {
      setErrors(response.graphQLErrors.map(error => error.message));
      if (!response.graphQLErrors.length) {
        localStorage.clear();
        authenticationExpiredCallback('')
        history.push('/sign-in', {authenticationExpired: true});
      }
    }
  });

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const updatePasswordHandler = () => {
    updatePassword({variables: {passwordId: passwordEntry.id, name: name, password: password}});
  };

  const deletePasswordHandler = () => {
    const confirm = window.confirm("Delete password?");
    if(confirm === true){
      deletePassword({variables: {passwordId: passwordEntry.id}});
    }
  };

  const errorMessage = errors ? errors.map((error, index) => {
    return <ErrorMessage key={index}>{error}</ErrorMessage>
  }) : null;

  if (deleted) return null;

  return (
    <div className='password-edit-container'>
      {errorMessage}
      <div className='password-edit'>
        <FormInput type='text' onChange={event => setName(event.target.value)} value={name} />
        <FormInput type={showPassword ? 'text' : 'password'} onChange={event => setPassword(event.target.value)} value={password} />
        <Button tooltip='Show password' onClick={() => togglePasswordVisibility()}>
          <FontAwesomeIcon icon={faEye} />
        </Button>
        <Button tooltip='Save edited password' onClick={() => updatePasswordHandler()} disabled={updateLoading ?? false}>
          <FontAwesomeIcon icon={faSave} />
        </Button>
        <Button tooltip='Delete password' onClick={() => deletePasswordHandler()} disabled={deleteLoading ?? false}>
          <FontAwesomeIcon icon={faTrash} />
        </Button>
      </div>
    </div>
  );
};

export default PasswordEdit;