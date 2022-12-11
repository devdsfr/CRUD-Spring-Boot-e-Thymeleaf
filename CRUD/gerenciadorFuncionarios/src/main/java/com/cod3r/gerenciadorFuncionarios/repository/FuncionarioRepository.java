package com.cod3r.gerenciadorFuncionarios.repository;

import com.cod3r.gerenciadorFuncionarios.model.Funcionario;
import com.cod3r.gerenciadorFuncionarios.model.FuncionarioSetor;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface FuncionarioRepository extends JpaRepository<Funcionario, Integer> {
    List<Funcionario> findBySetor(FuncionarioSetor funcionarioSetor);
}
