import {gql} from '@apollo/react-hooks';

export default gql`
  mutation CreatePassword($userId: ID!, $name: String!, $password: String!){
    createPassword(input: {userId: $userId, name: $name, password: $password}){
      id
      name
      password
    }
  }
`;