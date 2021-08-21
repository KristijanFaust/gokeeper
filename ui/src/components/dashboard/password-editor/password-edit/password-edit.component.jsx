import {faEye, faSave, faTrash} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {useState} from 'react';

import FormInput from '../../../input/input.component';
import Button from '../../../button/button.component';

import './password-edit.component.scss';

const PasswordEdit = ({passwordEntry}) => {
  const [name, setName] = useState(passwordEntry.name)
  const [password, setPassword] = useState(passwordEntry.password)
  const [showPassword, setShowPassword] = useState(false)

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <div className='password-edit'>
      <FormInput type='text' onChange={event => setName(event.target.value)} value={name} />
      <FormInput type={showPassword ? 'text' : 'password'} onChange={event => setPassword(event.target.value)} value={password} />
      <Button onClick={() => togglePasswordVisibility()}><FontAwesomeIcon icon={faEye} /></Button>
      <Button><FontAwesomeIcon icon={faSave} /></Button>
      <Button><FontAwesomeIcon icon={faTrash} /></Button>
    </div>
  );
};

export default PasswordEdit;