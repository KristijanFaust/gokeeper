import {gql} from '@apollo/react-hooks';

export default gql`
    query QueryUserPassword($userId: String!) {
        queryUserPasswords(userId: $userId) {
            id
            name
            password
        }
    }
`;