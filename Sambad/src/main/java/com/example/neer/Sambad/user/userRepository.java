package com.example.neer.Sambad.user;

import java.util.List;

import org.springframework.data.mongodb.repository.MongoRepository;

public interface userRepository extends MongoRepository<user,String>{
     List<user> findAllByStatus(Status status);

}
