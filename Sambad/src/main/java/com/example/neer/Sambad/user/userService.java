package com.example.neer.Sambad.user;

import java.util.List;

import org.springframework.stereotype.Service;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class userService {

    private final userRepository repository;

    public void saveUser(user user) {
        user.setStatus(Status.ONLINE);
        repository.save(user);
    }

    public void disconnect(user user) {
        var storedUser = repository.findById(user.getNickName()).orElse(null);
        if (storedUser != null) {
            storedUser.setStatus(Status.OFFLINE);
            repository.save(storedUser);
        }
    }

    public List<user> findConnectedUsers() {
        return repository.findAllByStatus(Status.ONLINE);
        
    }
}
