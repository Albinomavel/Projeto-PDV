using Avalonia.Controls;
using Avalonia.Interactivity;
using System;
using System.Net.Http;
using System.Threading.Tasks;
using System.Text.Json;
using System.Text.Json.Serialization;


namespace Client;

public partial class MainWindow : Window
{
    public MainWindow()
    {
        InitializeComponent();
    }
}