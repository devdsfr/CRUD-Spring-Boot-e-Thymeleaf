package com.cod3r.gerenciadorFuncionarios.contoller;

import com.cod3r.gerenciadorFuncionarios.model.Funcionario;
import com.cod3r.gerenciadorFuncionarios.model.FuncionarioSetor;
import com.cod3r.gerenciadorFuncionarios.repository.FuncionarioRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;

@Controller
public class HomeController {

    @Autowired
    FuncionarioRepository funcionarioRepository;

    @GetMapping("/home")
    public String home(Model model) {
//        return "home";
//    }

        List<Funcionario> funcionarios = funcionarioRepository.findAll();

        model.addAttribute("funcionarios", funcionarios);
        return "home";
    }

    @GetMapping
    public String setor(@RequestParam String setor, Model model) {
        FuncionarioSetor funcionarioSetor = FuncionarioSetor.valueOf(setor.toUpperCase());
        List<Funcionario> funcionarios = funcionarioRepository.findBySetor(funcionarioSetor);

        model.addAttribute("funcionarios", funcionarios);
        return "home";
    }
}
