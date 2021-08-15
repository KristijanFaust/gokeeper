import './input.styles.scss';

const FormInput = ({ label, ...otherProps }) => (
  <div className='group'>
    <input className='input' {...otherProps} />
    <label className='input-label'>{label}</label>
  </div>
);

export default FormInput;
